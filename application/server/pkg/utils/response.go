package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseWithData 返回带数据的成功响应
func ResponseWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// ResponseSuccess 返回成功响应（无数据）
func ResponseSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    nil,
	})
}

// ResponseError 返回错误响应
func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ResponseBadRequest 返回400错误
func ResponseBadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "请求参数错误"
	}
	ResponseError(c, http.StatusBadRequest, message)
}

// ResponseUnauthorized 返回401错误
func ResponseUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权访问"
	}
	ResponseError(c, http.StatusUnauthorized, message)
}

// ResponseForbidden 返回403错误
func ResponseForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "禁止访问"
	}
	ResponseError(c, http.StatusForbidden, message)
}

// ResponseNotFound 返回404错误
func ResponseNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	ResponseError(c, http.StatusNotFound, message)
}

// ResponseInternalServerError 返回500错误
func ResponseInternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	ResponseError(c, http.StatusInternalServerError, message)
}
