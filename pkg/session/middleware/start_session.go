package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/session"
)

func StartSession(cfg config.SessionConfig, storage session.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if storage == nil {
			panic("Session storage not configured. Please configure it first in bootstrap/web/web.go")
		}

		// Initialize session data
		var data *session.Data
		sessionId, err := ctx.Cookie(cfg.CookieName)

		if err == nil {
			// Get session data from storage
			sess, err := storage.Get(ctx, sessionId)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "unable_to_get_session_data",
					"data":    nil,
				})
				return
			}
			if sess != nil {
				data = sess
			}
		}
		if data == nil {
			data = session.NewEmptyData(cfg, ctx, storage)
			if err := data.Save(); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "unable_to_save_session_data",
					"data":    nil,
				})
				return
			}
		}
		ctx.Set("session", data)
		session.AddCookieToResponse(cfg, ctx, data.Id())
		ctx.Next()
	}
}
