package controller

import (
	blockDto "grets_server/dto/block_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// GlobalBlockController 全局区块控制器实例
var GlobalBlockController *BlockController

// InitBlockController 初始化区块控制器
func InitBlockController() {
	GlobalBlockController = NewBlockController()
}

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

// QueryBlockList 查询区块列表
func QueryBlockList(ctx *gin.Context) {
	GlobalBlockController.QueryBlockList(ctx)
}
