package controller

import (
	pictureDTO "grets_server/dto/picture_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// PictureController 图片控制器结构体
type PictureController struct {
	pictureService service.PictureService
}

// NewPictureController 创建图片控制器实例
func NewPictureController() *PictureController {
	return &PictureController{
		pictureService: service.NewPictureService(),
	}
}

// UploadPicture 上传图片
func (c *PictureController) UploadPicture(ctx *gin.Context) {
	var req pictureDTO.UploadPictureDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(ctx, err.Error())
		return
	}
	err := c.pictureService.UploadPicture(&req)
	if err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}
	utils.ResponseSuccess(ctx, "图片上传成功", nil)
}

// 创建全局图片控制器实例
var GlobalPictureController *PictureController

// 初始化图片控制器
func InitPictureController() {
	GlobalPictureController = NewPictureController()
}

// 为兼容现有路由，提供这些函数
func UploadPicture(c *gin.Context) {
	GlobalPictureController.UploadPicture(c)
}
