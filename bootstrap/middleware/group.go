package middleware

import (
	sessionsMiddleware "bitbucket.org/dptsi/go-framework/sessions/middleware"
	"bitbucket.org/dptsi/go-framework/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type MiddlewareGroup struct {
	i *do.Injector
}

func NewMiddlewareGroup(i *do.Injector) *MiddlewareGroup {
	return &MiddlewareGroup{
		i: i,
	}
}

func (m *MiddlewareGroup) GlobalMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		do.MustInvoke[*middleware.HandleCors](m.i).Execute,
	}
}

func (m *MiddlewareGroup) WebMiddleware() []gin.HandlerFunc {
	startSession := do.MustInvoke[*sessionsMiddleware.StartSession](m.i)
	verifyCsrfToken := do.MustInvoke[*sessionsMiddleware.VerifyCSRFToken](m.i)
	return []gin.HandlerFunc{
		startSession.Execute,
		verifyCsrfToken.Execute,
	}
}
