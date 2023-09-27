package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/session"
)

func StartSession(storage session.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := do.MustInvoke[config.Config](do.DefaultInjector).Session()

		// Initialize session data
		var data session.Data
		sessionId, err := ctx.Cookie(cfg.CookieName)

		if err != nil {
			// Generate new session id if not exist
			sessionId = uuid.NewString()
			data = session.NewData(ctx, sessionId, make(map[string]interface{}), storage)
		} else {
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
				data = *sess
			}
		}

		ctx.Set("session", data)

		// Set session cookie
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie(cfg.CookieName, sessionId, cfg.Lifetime, cfg.CookiePath, cfg.Domain, cfg.Secure, true)
		if err := data.Save(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "unable_to_save_session_data",
				"data":    nil,
			})
			return
		}
		ctx.Next()
	}
}
