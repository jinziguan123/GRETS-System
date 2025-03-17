package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseWithData 返回带数据的成功响应
func ResponseWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "操作成功",
		Data:    data,
	})
}

// ResponseSuccess 返回成功响应
func ResponseSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "操作成功",
	})
}

// ResponseBadRequest 返回400错误
func ResponseBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
	})
}

// ResponseUnauthorized 返回401错误
func ResponseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
	})
}

// ResponseForbidden 返回403错误
func ResponseForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
	})
}

// ResponseNotFound 返回404错误
func ResponseNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
	})
}

// ResponseInternalServerError 返回500错误
func ResponseInternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
	})
}
