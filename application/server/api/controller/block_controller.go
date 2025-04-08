package controller

import (
	blockDto "grets_server/dto/block_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

type BlockController struct {
	blockService service.BlockService
}

func NewBlockController() *BlockController {
	return &BlockController{
		blockService: service.NewBlockService(),
	}
}

func (c *BlockController) QueryBlockList(ctx *gin.Context) {
	// 绑定参数
	var queryBlockDTO blockDto.QueryBlockDTO
	if err := ctx.ShouldBindJSON(&queryBlockDTO); err != nil {
		utils.ResponseBadRequest(ctx, err.Error())
		return
	}

	block, err := c.blockService.QueryBlockList(queryBlockDTO)
	if err != nil {
		utils.ResponseInternalServerError(ctx, err.Error())
		return
	}
	utils.ResponseSuccess(ctx, "获取区块成功", block)
}
