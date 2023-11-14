package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func globalErrorHandler(isDebugMode bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		err := ctx.Errors.Last()
		if err == nil {
			return
		}
		requestId := ""
		reqIdInterface, exists := ctx.Get("request_id")
		if exists {
			if reqId, ok := reqIdInterface.(string); ok {
				requestId = reqId
			}
		}

		data := gin.H{
			"request_id": requestId,
		}
		log.Printf("Request ID: %s; Error: %s\n", requestId, err.Error())
		if isDebugMode {
			data["error"] = err.Error()
		}
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    9999,
				"message": "internal_server_error",
				"data":    data,
			},
		)
		ctx.Abort()
	}
}
