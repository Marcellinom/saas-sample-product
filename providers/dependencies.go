package bootstrap

import (
	"os"
	"strconv"

	"bitbucket.org/dptsi/go-framework/auth"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/sessions"

	"github.com/samber/do"
)

func CreateObjects(i *do.Injector) {
	// do.Provide[*middleware.HandleCors](i, func(i *do.Injector) (*middleware.HandleCors, error) {
	// 	return &middleware.HandleCors{
	// 		AllowedOrigins:   []string{"http://localhost:3000"},
	// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 		AllowedHeaders:   []string{"content-type", "x-csrf-token"},
	// 		ExposedHeaders:   []string{},
	// 		AllowCredentials: true,
	// 		MaxAge:           0,
	// 	}, nil
	// })
	// do.Provide[*appMiddleware.MiddlewareGroup](i, func(i *do.Injector) (*appMiddleware.MiddlewareGroup, error) {
	// 	return appMiddleware.NewMiddlewareGroup(i), nil
	// })
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
	do.Provide[contracts.SessionCookieWriter](i, func(i *do.Injector) (contracts.SessionCookieWriter, error) {
		return sessions.NewCookieUtil(sessionConfig), nil
	})
	do.Provide[contracts.AuthService](i, func(i *do.Injector) (contracts.AuthService, error) {
		driver := auth.NewSessionGuard(do.MustInvoke[contracts.SessionStorage](i), do.MustInvoke[contracts.SessionCookieWriter](i))

		return auth.NewService(
			auth.Config{
				Guards: map[string]auth.GuardsConfig{
					"web": {
						Driver: driver,
					},
				},
			},
		), nil
	})

	// do.Provide[*middleware.StartSession](i, func(i *do.Injector) (*middleware.StartSession, error) {
	// 	return middleware.NewStartSession(
	// 		sessionConfig,
	// 		do.MustInvoke[contracts.SessionStorage](i),
	// 		do.MustInvoke[contracts.SessionCookieWriter](i),
	// 	), nil
	// })
	// do.Provide[*middleware.VerifyCSRFToken](i, func(i *do.Injector) (*middleware.VerifyCSRFToken, error) {
	// 	return middleware.NewVerifyCSRFToken(), nil
	// })
	// do.Provide[*middleware.ActiveRole](i, func(i *do.Injector) (*middleware.ActiveRole, error) {
	// 	return middleware.NewActiveRole(do.MustInvoke[contracts.AuthService](i)), nil
	// })
}
