package oidc

import (
	"context"
	"errors"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"its.ac.id/base-go/pkg/session"
)

var (
	ErrNoEndSessionEndpoint = errors.New("this_oidc_provider_does_not_support_end_session_endpoint")
	ErrInvalidState         = errors.New("invalid_state")
	ErrInvalidNonce         = errors.New("invalid_nonce")
	ErrInvalidCodeChallenge = errors.New("invalid_code_challenge")
	ErrInvalidIdToken       = errors.New("invalid_id_token")
	ErrRetrieveUserInfo     = errors.New("error_retrieve_user_info")
)

const (
	stateKey        = "oidc.state"
	idTokenKey      = "oidc.id_token"
	nonceKey        = "oidc.nonce"
	codeVerifierKey = "oidc.code_verifier"
)

type Client struct {
	provider    *oidc.Provider
	oauthConfig oauth2.Config
	sess        *session.Data

	needToVerifyState bool
	needToVerifyNonce bool
	isPKCEEnabled     bool
}

func NewClient(
	ctx context.Context,
	providerUrl string,
	clientID string,
	clientSecret string,
	redirectURL string,
	scopes []string,
	sess *session.Data,
) (*Client, error) {
	provider, err := oidc.NewProvider(ctx, providerUrl)
	if err != nil {
		return nil, err
	}
	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		Scopes: scopes,
	}

	return &Client{provider, cfg, sess, true, true, true}, nil
}

func (c *Client) SetNeedToVerifyState(needToVerifyState bool) {
	c.needToVerifyState = needToVerifyState
}

func (c *Client) SetNeedToVerifyNonce(needToVerifyNonce bool) {
	c.needToVerifyNonce = needToVerifyNonce
}

func (c *Client) SetPKCEEnabled(isPKCEEnabled bool) {
	c.isPKCEEnabled = isPKCEEnabled
}

func (c *Client) RedirectURL() (string, error) {
	authCodeOptions := make([]oauth2.AuthCodeOption, 0)

	if c.isPKCEEnabled {
		codeVerifier := oauth2.GenerateVerifier()
		c.sess.Set(codeVerifierKey, codeVerifier)
		authCodeOptions = append(authCodeOptions, oauth2.S256ChallengeOption(codeVerifier))
	}

	if c.needToVerifyNonce {
		nonce := uuid.NewString()
		c.sess.Set(nonceKey, nonce)
		authCodeOptions = append(authCodeOptions, oauth2.SetAuthURLParam("nonce", nonce))
	}

	state := ""
	if c.needToVerifyState {
		state = uuid.NewString()
		c.sess.Set(stateKey, state)
	}

	isNeedToSaveSession := c.needToVerifyState || c.needToVerifyNonce || c.isPKCEEnabled
	if err := c.sess.Save(); isNeedToSaveSession && err != nil {
		return "", err
	}

	return c.oauthConfig.AuthCodeURL(
		state,
		authCodeOptions...,
	), nil
}

func (c *Client) ExchangeCodeForToken(ctx context.Context, code string, state string) (*oauth2.Token, *oidc.IDToken, error) {
	if err := c.verifyState(state); c.needToVerifyState && err != nil {
		return nil, nil, err
	}

	authCodeOptions := make([]oauth2.AuthCodeOption, 0)
	codeVerifier, err := c.GetCodeVerifierAndRemoveFromSession()
	if c.isPKCEEnabled && err != nil {
		return nil, nil, err
	}
	if c.isPKCEEnabled {
		authCodeOptions = append(authCodeOptions, oauth2.VerifierOption(codeVerifier))
	}

	token, err := c.oauthConfig.Exchange(ctx, code, authCodeOptions...)
	if err != nil {
		return nil, nil, err
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, nil, errors.New("no_id_token_in_payload")
	}

	IDToken, err := c.parseAndVerifyIDToken(ctx, rawIDToken)
	if err != nil {
		return nil, nil, err
	}

	if c.needToVerifyNonce {
		c.verifyNonce(IDToken.Nonce)
	}

	c.sess.Set(idTokenKey, rawIDToken)
	if err := c.sess.Save(); err != nil {
		return nil, nil, err
	}

	return token, IDToken, nil
}

func (c *Client) parseAndVerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	// Parse and verify ID Token payload.
	var verifier = c.provider.Verifier(&oidc.Config{ClientID: c.oauthConfig.ClientID})
	parsed, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, ErrInvalidIdToken
	}

	return parsed, nil
}

func (c *Client) UserInfo(ctx context.Context, t *oauth2.Token) (*oidc.UserInfo, error) {
	userInfo, err := c.provider.UserInfo(ctx, oauth2.StaticTokenSource(t))
	if err != nil {
		return nil, ErrRetrieveUserInfo
	}

	return userInfo, nil
}

func (c *Client) RPInitiatedLogout(postLogoutRedirectURI string) (string, error) {
	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	if err := c.provider.Claims(&claims); err != nil {
		return "", err
	}
	endSessionEndpoint := claims.EndSessionEndpoint
	if endSessionEndpoint == "" {
		return "", ErrNoEndSessionEndpoint
	}
	req, err := http.NewRequest("GET", endSessionEndpoint, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()

	idTokenHintItf, exists := c.sess.Get(idTokenKey)
	if idTokenHint, ok := idTokenHintItf.(string); exists && ok && idTokenHint != "" {
		q.Add("id_token_hint", idTokenHint)
	}
	c.sess.Delete(idTokenKey)
	c.sess.Save()

	if postLogoutRedirectURI != "" {
		q.Add("post_logout_redirect_uri", postLogoutRedirectURI)
	}

	req.URL.RawQuery = q.Encode()
	return req.URL.String(), nil
}

func (c *Client) verifyState(state string) error {
	cookieState, ok := c.sess.Get(stateKey)
	if !ok {
		cookieState = ""
	}
	c.sess.Delete(stateKey)
	c.sess.Save()

	if state == "" || state != cookieState {
		return ErrInvalidState
	}

	return nil
}

func (c *Client) GetCodeVerifierAndRemoveFromSession() (string, error) {
	codeVerifierIf, ok := c.sess.Get(codeVerifierKey)
	c.sess.Delete(codeVerifierKey)
	if err := c.sess.Save(); err != nil {
		return "", err
	}
	codeVerifier := ""
	if ok {
		codeVerifier, ok = codeVerifierIf.(string)
		if !ok {
			codeVerifier = ""
		}
	}

	if codeVerifier == "" {
		return "", ErrInvalidCodeChallenge
	}

	return codeVerifier, nil
}

func (c *Client) verifyNonce(nonce string) error {
	cookieNonce, ok := c.sess.Get(nonceKey)
	if !ok {
		cookieNonce = ""
	}
	c.sess.Delete(nonceKey)
	c.sess.Save()

	if nonce != cookieNonce {
		return ErrInvalidNonce
	}

	return nil
}
