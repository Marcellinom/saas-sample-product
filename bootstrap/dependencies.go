package bootstrap

import (
	"os"
	"strconv"

	"bitbucket.org/dptsi/go-framework/auth"
	authMiddleware "bitbucket.org/dptsi/go-framework/auth/middleware"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/database"
	"bitbucket.org/dptsi/go-framework/sessions"
	sessionsMiddleware "bitbucket.org/dptsi/go-framework/sessions/middleware"
	webMiddleware "bitbucket.org/dptsi/go-framework/web/middleware"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/middleware"
)

func CreateObjects(i *do.Injector) {
	do.Provide[*webMiddleware.HandleCors](i, func(i *do.Injector) (*webMiddleware.HandleCors, error) {
		return &webMiddleware.HandleCors{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"content-type", "x-csrf-token"},
			ExposedHeaders:   []string{},
			AllowCredentials: true,
			MaxAge:           0,
		}, nil
	})
	do.Provide[*middleware.MiddlewareGroup](i, func(i *do.Injector) (*middleware.MiddlewareGroup, error) {
		return middleware.NewMiddlewareGroup(i), nil
	})
	do.Provide[*database.Manager](i, func(i *do.Injector) (*database.Manager, error) {
		return database.NewManager(), nil
	})
	sessionMaxAge, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	if err != nil {
		sessionMaxAge = 86400
	}
	sessionConfig := sessions.SessionsConfig{
		Name:           os.Getenv("SESSION_NAME"),
		CsrfCookieName: os.Getenv("SESSION_CSRF_COOKIE_NAME"),
		MaxAge:         sessionMaxAge,
		Path:           os.Getenv("SESSION_PATH"),
		Domain:         os.Getenv("SESSION_DOMAIN"),
		Secure:         os.Getenv("SESSION_SECURE") == "true",
	}
	do.Provide[*sessions.CookieUtil](i, func(i *do.Injector) (*sessions.CookieUtil, error) {
		return sessions.NewCookieUtil(sessionConfig), nil
	})

	do.Provide[*sessionsMiddleware.StartSession](i, func(i *do.Injector) (*sessionsMiddleware.StartSession, error) {
		return sessionsMiddleware.NewStartSession(
			sessionConfig,
			do.MustInvoke[contracts.SessionStorage](i),
			*(do.MustInvoke[*sessions.CookieUtil](i)),
		), nil
	})
	do.Provide[*sessionsMiddleware.VerifyCSRFToken](i, func(i *do.Injector) (*sessionsMiddleware.VerifyCSRFToken, error) {
		return sessionsMiddleware.NewVerifyCSRFToken(), nil
	})
	do.Provide[*auth.Service](i, func(i *do.Injector) (*auth.Service, error) {
		return auth.NewService(
			do.MustInvoke[contracts.SessionStorage](i),
		), nil
	})
	do.Provide[*authMiddleware.ActiveRole](i, func(i *do.Injector) (*authMiddleware.ActiveRole, error) {
		return authMiddleware.NewActiveRole(*do.MustInvoke[*auth.Service](i)), nil
	})
}
