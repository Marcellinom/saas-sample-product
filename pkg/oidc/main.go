package oidc

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	StateKey    = "oidc_state"
	IdTokenKey  = "oidc_id_token"
	StateMaxAge = 60 * 5 // 5 minutes

	AuthorizationCodeNotFound = "authorization_code_not_found"
	InvalidState              = "invalid_state"
	InvalidIdToken            = "invalid_id_token"
	ErrorRetrieveUserInfo     = "error_retrieve_user_info"
)

type CookieProvider interface {
	Cookie(name string) (string, error)
	SetCookie(name string, value string)
}

type QueryParamsProvider interface {
	GetQuery(key string) (string, bool)
}

type Client struct {
	p           *oidc.Provider
	cp          CookieProvider
	qp          QueryParamsProvider
	ctx         context.Context
	verifyState bool
}

func NewClient(ctx context.Context, pUrl string, cp CookieProvider, qp QueryParamsProvider) (*Client, error) {
	provider, err := oidc.NewProvider(ctx, pUrl)
	if err != nil {
		return nil, err
	}

	return &Client{provider, cp, qp, ctx, true}, nil
}

func (c *Client) SetVerifyState(verifyState bool) {
	c.verifyState = verifyState
}

func (c *Client) RedirectURL(clientID string, clientSecret string, redirectURL string, scopes []string) string {
	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: c.p.Endpoint(),

		Scopes: scopes,
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	state := uuid.NewString()
	c.cp.SetCookie(StateKey, state)
	return cfg.AuthCodeURL(state)
}

func (c *Client) UserInfo(clientID string, clientSecret string, redirectURL string, scopes []string) (*oidc.UserInfo, error) {
	code, exist := c.qp.GetQuery("code")
	if !exist {
		return nil, errors.New(AuthorizationCodeNotFound)
	}
	if c.verifyState {
		state, exist := c.qp.GetQuery("state")
		cookieState, err := c.cp.Cookie(StateKey)
		c.cp.SetCookie(StateKey, "")
		if err != nil {
			return nil, err
		}

		if state == "" || !exist || state != cookieState {
			return nil, errors.New(InvalidState)
		}
	}

	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: c.p.Endpoint(),

		Scopes: scopes,
	}

	token, err := cfg.Exchange(c.ctx, code)
	if err != nil {
		return nil, err
	}
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New(InvalidIdToken)
	}

	// Parse and verify ID Token payload.
	var verifier = c.p.Verifier(&oidc.Config{ClientID: clientID})
	_, err = verifier.Verify(c.ctx, rawIDToken)
	if err != nil {
		return nil, errors.New(InvalidIdToken)
	}
	c.cp.SetCookie(IdTokenKey, rawIDToken)
	userInfo, err := c.p.UserInfo(c.ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, errors.New(ErrorRetrieveUserInfo)
	}

	return userInfo, nil
}
