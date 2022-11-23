package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpOK(ctx *gin.Context, body gin.H) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"Error": nil,
		"Data":  body,
	})
}

func HttpError(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(code, gin.H{
		"Error": err.Error(),
		"Data":  nil,
	})
}
