package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/app/common/errors"
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

func AbortAndResponseErrorWithJSON(c *gin.Context, err error) {
	if notFound, isNotFound := err.(*errors.NotFoundError); isNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": notFound.Error(),
			"data":    nil,
		})
		return
	}
	if mismatch, isVersionMismatch := err.(*errors.AggregateVersionMismatchError); isVersionMismatch {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": mismatch.Error(),
			"data":    nil,
		})
		return
	}
	if invariantError, isInvariantError := err.(*errors.InvariantError); isInvariantError {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": invariantError.Error(),
			"data":    nil,
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, InternalServerErrorResponse)
}

func HandleInfiniteScrollResponse[T any](ctx *gin.Context, limit int, result *InfiniteScrollResult[T], err error) {
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"links": map[string]string{
			"next": fmt.Sprintf("%s://%s?cursor=%s&limit=%d", scheme, ctx.Request.Host+ctx.Request.URL.Path, result.NextCursor(), limit),
		},
		"meta": map[string]interface{}{
			"total": result.Total(),
		},
		"data": result.Data(),
	})
}

func HandleTableAdvancedResponse[T any](ctx *gin.Context, limit int, currentPage int, result *TableAdvancedResult[T], err error) {
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"links": map[string]string{
			"prev": fmt.Sprintf("%s://%s?end_before=%s&page=%d", scheme, ctx.Request.Host+ctx.Request.URL.Path, result.EndBefore(), max(1, currentPage-1)),
			"next": fmt.Sprintf("%s://%s?start_after=%s&page=%d", scheme, ctx.Request.Host+ctx.Request.URL.Path, result.StartAfter(), min(result.Total()/limit+1, currentPage+1)),
		},
		"meta": map[string]interface{}{
			"total": result.Total(),
			"range": result.ItemCount(),
			"page":  currentPage,
		},
		"data": result.Data(),
	})
}
