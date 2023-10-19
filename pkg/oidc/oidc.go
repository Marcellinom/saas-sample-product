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
	ErrNoEndSessionEndpoint = errors.New("no end session endpoint configured. Please add OIDC_END_SESSION_ENDPOINT to your .env file")
)

const (
	stateKey   = "oidc.state"
	idTokenKey = "oidc.id_token"
	nonceKey   = "oidc.nonce"

	AuthorizationCodeNotFound = "authorization_code_not_found"
	InvalidState              = "invalid_state"
	InvalidNonce              = "invalid_nonce"
	InvalidIdToken            = "invalid_id_token"
	ErrorRetrieveUserInfo     = "error_retrieve_user_info"
)

type Client struct {
	provider    *oidc.Provider
	oauthConfig oauth2.Config
	verifyState bool
	sess        *session.Data
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

	return &Client{provider, cfg, true, sess}, nil
}

func (c *Client) SetVerifyState(verifyState bool) {
	c.verifyState = verifyState
}

func (c *Client) RedirectURL() string {
	state := uuid.NewString()
	nonce := uuid.NewString()
	c.sess.Set(stateKey, state)
	c.sess.Set(nonceKey, nonce)
	c.sess.Save()
	return c.oauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
	)
}

func (c *Client) ExchangeCodeForToken(ctx context.Context, code string, state string) (*oauth2.Token, *oidc.IDToken, error) {
	if c.verifyState {
		cookieState, ok := c.sess.Get(stateKey)
		if !ok {
			cookieState = ""
		}
		c.sess.Delete(stateKey)
		c.sess.Save()

		if state == "" || state != cookieState {
			return nil, nil, errors.New(InvalidState)
		}
	}

	token, err := c.oauthConfig.Exchange(ctx, code)
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

	nonceIf, ok := c.sess.Get(nonceKey)
	nonce := ""
	if ok {
		nonce, ok = nonceIf.(string)
		if !ok {
			nonce = ""
		}
	}

	if nonce != "" && IDToken.Nonce != nonce {
		return nil, nil, errors.New(InvalidNonce)
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
		return nil, errors.New(InvalidIdToken)
	}

	return parsed, nil
}

func (c *Client) UserInfo(ctx context.Context, t *oauth2.Token) (*oidc.UserInfo, error) {
	userInfo, err := c.provider.UserInfo(ctx, oauth2.StaticTokenSource(t))
	if err != nil {
		return nil, errors.New(ErrorRetrieveUserInfo)
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
