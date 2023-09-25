package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var UnauthorizedResponse = gin.H{
	"code":    http.StatusUnauthorized,
	"message": "unauthorized",
	"data":    nil,
}

var InternalServerErrorResponse = gin.H{
	"code":    http.StatusInternalServerError,
	"message": "internal_server_error",
	"data":    nil,
}
