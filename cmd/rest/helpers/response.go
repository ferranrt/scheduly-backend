package helpers

import "github.com/gin-gonic/gin"

func BuildErrorResponse(err string) gin.H {
	return gin.H{
		"error": err,
	}
}

func BuildSuccessResponse(data interface{}) gin.H {
	return gin.H{
		"data": data,
	}
}
