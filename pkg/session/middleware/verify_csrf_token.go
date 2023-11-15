package middleware

import (
	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/app/common/errors"
	"its.ac.id/base-go/pkg/session"
)

var errInvalidCSRFToken = errors.NewForbiddenError("invalid_csrf_token", "")
var methodsWithoutCSRFToken = []string{"GET", "HEAD", "OPTIONS"}

func VerifyCSRFToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := session.Default(ctx)
		sessionCSRFToken := sess.CSRFToken()
		requestCSRFToken := ctx.Request.Header.Get("X-CSRF-TOKEN")

		// Skip CSRF token verification for some methods
		for _, method := range methodsWithoutCSRFToken {
			if ctx.Request.Method == method {
				ctx.Next()
				return
			}
		}

		if sessionCSRFToken == "" || sessionCSRFToken != requestCSRFToken {
			ctx.Error(errInvalidCSRFToken)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
