package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	commonErrors "its.ac.id/base-go/pkg/app/common/errors"
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

		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			errorData := commonErrors.GetValidationErrors(validationError)
			data["errors"] = errorData
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    9998,
					"message": "validation_error",
					"data":    data,
				},
			)
		} else {
			log.Printf("Request ID: %s; Status: 500; Error: %s\n", requestId, err.Error())
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
		}

		ctx.Abort()
	}
}
