package middleware

import (
	"bitbucket.org/dptsi/go-framework/http/middleware"
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
	startSession := do.MustInvoke[*middleware.StartSession](m.i)
	verifyCsrfToken := do.MustInvoke[*middleware.VerifyCSRFToken](m.i)
	return []gin.HandlerFunc{
		startSession.Execute,
		verifyCsrfToken.Execute,
	}
}
