package utils

import (
	"ginchat/common"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"code":    common.SUCCESS,
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(200, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
