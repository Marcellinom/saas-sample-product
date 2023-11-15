package middleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/session"
)

func StartSession(cfg config.SessionConfig, storage session.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if storage == nil {
			err := errors.New("session storage not configured. please configure it first in bootstrap/web/web.go")
			ctx.Error(fmt.Errorf("start session middleware: %w", err))
			ctx.Abort()
		}

		// Initialize session data
		var data *session.Data
		sessionId, err := ctx.Cookie(cfg.CookieName)

		if err == nil {
			// Get session data from storage
			sess, err := storage.Get(ctx, sessionId)
			if err != nil {
				ctx.Error(err)
				ctx.Abort()
				return
			}
			if sess != nil {
				data = sess
			}
		}
		if data == nil {
			data = session.NewEmptyData(cfg, ctx, storage)
			if err := data.Save(); err != nil {
				ctx.Error(fmt.Errorf("start session middleware: %w", err))
				ctx.Abort()
			}
		}
		ctx.Set("session", data)
		session.AddCookieToResponse(cfg, ctx, data.Id())
		ctx.Next()
	}
}
