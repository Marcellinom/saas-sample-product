package middleware

import (
	"bitbucket.org/dptsi/base-go-libraries/web/middleware"
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
		do.MustInvokeNamed[*middleware.HandleCors](m.i, "HandleCorsMiddleware").Execute,
	}
}

func (m *MiddlewareGroup) WebMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
