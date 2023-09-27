package session

import "github.com/gin-gonic/gin"

type Storage interface {
	Get(ctx *gin.Context, sessionId string) (*Data, error)
	Save(ctx *gin.Context, sessionId string, data map[string]interface{}) error
	Delete(ctx *gin.Context, sessionId string) error
}
