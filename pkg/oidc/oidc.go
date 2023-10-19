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

	verifyState bool
	verifyNonce bool
	enablePKCE  bool
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

func (c *Client) SetVerifyState(verifyState bool) {
	c.verifyState = verifyState
}

func (c *Client) RedirectURL() string {
	state := uuid.NewString()
	nonce := uuid.NewString()
	codeVerifier := oauth2.GenerateVerifier()

	c.sess.Set(codeVerifierKey, codeVerifier)
	c.sess.Set(stateKey, state)
	c.sess.Set(nonceKey, nonce)
	c.sess.Save()

	return c.oauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.S256ChallengeOption(codeVerifier),
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
			return nil, nil, ErrInvalidState
		}
	}

	var token *oauth2.Token
	if c.enablePKCE {
		codeVerifierIf, ok := c.sess.Get(codeVerifierKey)
		c.sess.Delete(codeVerifierKey)
		if err := c.sess.Save(); err != nil {
			return nil, nil, err
		}
		codeVerifier := ""
		if ok {
			codeVerifier, ok = codeVerifierIf.(string)
			if !ok {
				codeVerifier = ""
			}
		}

		if codeVerifier == "" {
			return nil, nil, ErrInvalidCodeChallenge
		}

		if tmpToken, err := c.oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier)); err != nil {
			return nil, nil, err
		} else {
			token = tmpToken
		}
	} else {
		if tmpToken, err := c.oauthConfig.Exchange(ctx, code); err != nil {
			return nil, nil, err
		} else {
			token = tmpToken
		}
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, nil, errors.New("no_id_token_in_payload")
	}

	IDToken, err := c.parseAndVerifyIDToken(ctx, rawIDToken)
	if err != nil {
		return nil, nil, err
	}

	if c.verifyNonce {
		nonceIf, ok := c.sess.Get(nonceKey)
		c.sess.Delete(nonceKey)
		if err := c.sess.Save(); err != nil {
			return nil, nil, err
		}
		nonce := ""
		if ok {
			nonce, ok = nonceIf.(string)
			if !ok {
				nonce = ""
			}
		}

		if nonce != "" && IDToken.Nonce != nonce {
			return nil, nil, ErrInvalidNonce
		}
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
