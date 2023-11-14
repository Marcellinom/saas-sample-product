package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/app/common/errors"
)

type TableAdvancedLinks struct {
	Next string `json:"next" example:"http://localhost:8080?start_after=2021-08-01&page=3"`
	Prev string `json:"prev" example:"http://localhost:8080?end_before=2021-08-01&page=1"`
}

type TableAdvancedMeta struct {
	Total int `json:"total" example:"100"`
	Range int `json:"range" example:"10"`
	Page  int `json:"page" example:"2"`
}

type TableAdvancedResponse[T any] struct {
	Code    int                `json:"code" example:"123"`
	Message string             `json:"message"`
	Links   TableAdvancedLinks `json:"links"`
	Meta    TableAdvancedMeta  `json:"meta"`
	Data    []T                `json:"data"`
}

type InfiniteScrollLinks struct {
	Next string `json:"next" example:"http://localhost:8080?cursor=2021-08-01&limit=10"`
}

type InfiniteScrollMeta struct {
	Total int `json:"total" example:"100"`
}

type InfiniteScrollResponse[T any] struct {
	Code    int                 `json:"code" example:"123"`
	Message string              `json:"message"`
	Links   InfiniteScrollLinks `json:"links"`
	Meta    InfiniteScrollMeta  `json:"meta"`
	Data    []T                 `json:"data"`
}

// DEPRECATED: Jangan dipakai untuk kode baru (alternatif menyusul)
var UnauthorizedResponse = gin.H{
	"code":    http.StatusUnauthorized,
	"message": "unauthorized",
	"data":    nil,
}

// DEPRECATED: Jangan dipakai untuk kode baru (alternatif menyusul)
var InternalServerErrorResponse = gin.H{
	"code":    http.StatusInternalServerError,
	"message": "internal_server_error",
	"data":    nil,
}

// DEPRECATED: Jangan dipakai untuk kode baru (alternatif menyusul)
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

	ctx.JSON(http.StatusOK, TableAdvancedResponse[T]{
		Code:    http.StatusOK,
		Message: "success",
		Links: TableAdvancedLinks{
			Prev: fmt.Sprintf("%s://%s?end_before=%s&page=%d", scheme, ctx.Request.Host+ctx.Request.URL.Path, result.EndBefore(), max(1, currentPage-1)),
			Next: fmt.Sprintf("%s://%s?start_after=%s&page=%d", scheme, ctx.Request.Host+ctx.Request.URL.Path, result.StartAfter(), min(result.Total()/limit+1, currentPage+1)),
		},
		Meta: TableAdvancedMeta{
			Total: result.Total(),
			Range: result.ItemCount(),
			Page:  currentPage,
		},
	})
}
