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

		var validationErrors validator.ValidationErrors
		var badRequestError commonErrors.BadRequestError
		if errors.As(err, &validationErrors) {
			errorData := commonErrors.GetValidationErrors(validationErrors)
			data["errors"] = errorData
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    statusCode[validationError],
					"message": validationError,
					"data":    data,
				},
			)
		} else if errors.As(err, &badRequestError) {
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
			for key, val := range badRequestError.Data() {
				data[key] = val
			}
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    badRequestError.Code(),
					"message": badRequestError.Message(),
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
					"code":    statusCode[internalServerError],
					"message": internalServerError,
					"data":    data,
				},
			)
		}

		ctx.Abort()
	}
}
