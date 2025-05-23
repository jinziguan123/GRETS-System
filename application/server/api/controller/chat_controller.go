package controller

import (
	chatDto "grets_server/dto/chat_dto"
	"grets_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var GlobalChatController *ChatController

func InitChatController() {
	GlobalChatController = NewChatController(service.GlobalChatService)
}

type ChatController struct {
	chatService service.ChatService
}

func NewChatController(chatService service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

// VerifyCapital 验资
// @Summary 验资
// @Description 验证用户资金是否充足以创建聊天室
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chatDto.VerifyCapitalDTO true "验资请求"
// @Success 200 {object} response.Response{data=chat_dto.VerifyCapitalResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/verify-capital [post]
func (ctrl *ChatController) VerifyCapital(c *gin.Context) {
	// 绑定请求参数
	var req chatDto.VerifyCapitalDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 调用服务层
	result, err := ctrl.chatService.VerifyCapital(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "验资完成",
		"data":    result,
	})
}

// CreateChatRoom 创建聊天室
// @Summary 创建聊天室
// @Description 验资后创建一对一聊天室
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat_dto.CreateChatRoomDTO true "创建聊天室请求"
// @Success 200 {object} response.Response{data=chat_dto.ChatRoomDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/room [post]
func (ctrl *ChatController) CreateChatRoom(c *gin.Context) {
	// 绑定请求参数
	var req chatDto.CreateChatRoomDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 调用服务层
	result, err := ctrl.chatService.CreateChatRoom(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "聊天室创建成功",
		"data":    result,
	})
}

// GetChatRoomList 获取聊天室列表
// @Summary 获取聊天室列表
// @Description 获取用户的聊天室列表
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chatDto.QueryChatRoomListDTO true "查询聊天室列表请求"
// @Success 200 {object} response.Response{data=chat_dto.ChatRoomListResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/getChatRoomList [post]
func (ctrl *ChatController) GetChatRoomList(c *gin.Context) {
	// 绑定请求参数
	var req chatDto.QueryChatRoomListDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}

	// 调用服务层
	result, err := ctrl.chatService.GetChatRoomList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取聊天室列表成功",
		"data":    result,
	})
}

// SendMessage 发送消息
// @Summary 发送消息
// @Description 在聊天室中发送消息
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat_dto.SendMessageDTO true "发送消息请求"
// @Success 200 {object} response.Response{data=chat_dto.ChatMessageDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/message [post]
func (ctrl *ChatController) SendMessage(c *gin.Context) {
	// 绑定请求参数
	var req chatDto.SendMessageDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 调用服务层
	result, err := ctrl.chatService.SendMessage(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "消息发送成功",
		"data":    result,
	})
}

// GetChatMessageList 获取聊天消息列表
// @Summary 获取聊天消息列表
// @Description 获取指定聊天室的消息列表
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chatDto.QueryChatMessageListDTO true "查询消息列表请求"
// @Success 200 {object} response.Response{data=chat_dto.ChatMessageListResponseDTO}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/getChatMessageList [post]
func (ctrl *ChatController) GetChatMessageList(c *gin.Context) {
	// 绑定请求参数
	var req chatDto.QueryChatMessageListDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}

	// 调用服务层
	result, err := ctrl.chatService.GetChatMessageList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取消息列表成功",
		"data":    result,
	})
}

// MarkMessagesRead 标记消息为已读
// @Summary 标记消息为已读
// @Description 标记指定聊天室的消息为已读
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat_dto.MarkMessagesReadDTO true "标记已读请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/message/read [put]
func (ctrl *ChatController) MarkMessagesRead(c *gin.Context) {
	// 绑定请求参数
	var req chatDto.MarkMessagesReadDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 调用服务层
	err := ctrl.chatService.MarkMessagesRead(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "标记已读成功",
	})
}

// CloseChatRoom 关闭聊天室
// @Summary 关闭聊天室
// @Description 关闭指定的聊天室（仅卖方可操作）
// @Tags 聊天室
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomUUID path string true "聊天室UUID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/chat/room/{roomUUID}/close [put]
func (ctrl *ChatController) CloseChatRoom(c *gin.Context) {
	var req chatDto.CloseChatRoomDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数绑定失败: " + err.Error(),
		})
		return
	}

	// 调用服务层
	err := ctrl.chatService.CloseChatRoom(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "聊天室关闭成功",
	})
}

func VerifyCapital(c *gin.Context) {
	GlobalChatController.VerifyCapital(c)
}

func CreateChatRoom(c *gin.Context) {
	GlobalChatController.CreateChatRoom(c)
}

func GetChatRoomList(c *gin.Context) {
	GlobalChatController.GetChatRoomList(c)
}

func SendMessage(c *gin.Context) {
	GlobalChatController.SendMessage(c)
}

func GetChatMessageList(c *gin.Context) {
	GlobalChatController.GetChatMessageList(c)
}

func MarkMessagesRead(c *gin.Context) {
	GlobalChatController.MarkMessagesRead(c)
}

func CloseChatRoom(c *gin.Context) {
	GlobalChatController.CloseChatRoom(c)
}
