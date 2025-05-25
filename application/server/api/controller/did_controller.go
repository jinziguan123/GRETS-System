package controller

import (
	"grets_server/constants"
	didDto "grets_server/dto/did_dto"
	"grets_server/pkg/utils"
	"grets_server/service"

	"github.com/gin-gonic/gin"
)

// DIDController DID控制器
type DIDController struct {
	didService service.DIDService
}

// NewDIDController 创建DID控制器
func NewDIDController(didService service.DIDService) *DIDController {
	return &DIDController{
		didService: didService,
	}
}

// CreateDID 创建DID
func (c *DIDController) CreateDID(ctx *gin.Context) {
	var req didDto.CreateDIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.CreateDID(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "DID创建成功", response)
}

// ResolveDID 解析DID
func (c *DIDController) ResolveDID(ctx *gin.Context) {
	did := ctx.Param("did")
	if did == "" {
		utils.ResponseError(ctx, constants.ParamError, "DID参数不能为空")
		return
	}

	response, err := c.didService.ResolveDID(did)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "DID解析成功", response)
}

// UpdateDIDDocument 更新DID文档
func (c *DIDController) UpdateDIDDocument(ctx *gin.Context) {
	var req didDto.UpdateDIDDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	if err := c.didService.UpdateDIDDocument(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "DID文档更新成功", nil)
}

// GetChallenge 获取认证挑战
func (c *DIDController) GetChallenge(ctx *gin.Context) {
	var req didDto.GetChallengeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.GetChallenge(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "获取认证挑战成功", response)
}

// DIDLogin DID登录
func (c *DIDController) DIDLogin(ctx *gin.Context) {
	var req didDto.DIDLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.DIDLogin(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "DID登录成功", response)
}

// IssueCredential 签发凭证
func (c *DIDController) IssueCredential(ctx *gin.Context) {
	var req didDto.IssueCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.IssueCredential(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "凭证签发成功", response)
}

// GetCredentials 获取凭证
func (c *DIDController) GetCredentials(ctx *gin.Context) {
	var req didDto.GetCredentialsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.GetCredentials(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "获取凭证成功", response)
}

// VerifyPresentation 验证展示
func (c *DIDController) VerifyPresentation(ctx *gin.Context) {
	var req didDto.VerifyPresentationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.VerifyPresentation(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "验证展示成功", response)
}

// RevokeCredential 撤销凭证
func (c *DIDController) RevokeCredential(ctx *gin.Context) {
	var req didDto.RevokeCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	if err := c.didService.RevokeCredential(&req); err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "凭证撤销成功", nil)
}

// DIDRegister DID注册（兼容传统注册）
func (c *DIDController) DIDRegister(ctx *gin.Context) {
	var req didDto.DIDRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, constants.ParamError, "参数错误: "+err.Error())
		return
	}

	response, err := c.didService.DIDRegister(&req)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, "DID注册成功", response)
}

// GetDIDByUser 根据用户信息获取DID
func (c *DIDController) GetDIDByUser(ctx *gin.Context) {
	citizenID := ctx.Query("citizenID")
	organization := ctx.Query("organization")

	if citizenID == "" || organization == "" {
		utils.ResponseError(ctx, constants.ParamError, "身份证号和组织不能为空")
		return
	}

	did, err := c.didService.GetDIDByUser(citizenID, organization)
	if err != nil {
		utils.ResponseError(ctx, constants.ServiceError, err.Error())
		return
	}

	if did == "" {
		utils.ResponseError(ctx, constants.ServiceError, "用户DID不存在")
		return
	}

	utils.ResponseSuccess(ctx, "获取用户DID成功", gin.H{"did": did})
}

// 创建全局DID控制器实例
var GlobalDIDController *DIDController

// 初始化DID控制器
func InitDIDController() {
	GlobalDIDController = NewDIDController(service.GlobalDIDService)
}

// 为兼容现有路由，提供这些函数
func CreateDID(c *gin.Context) {
	GlobalDIDController.CreateDID(c)
}

func ResolveDID(c *gin.Context) {
	GlobalDIDController.ResolveDID(c)
}

func UpdateDIDDocument(c *gin.Context) {
	GlobalDIDController.UpdateDIDDocument(c)
}

func GetChallenge(c *gin.Context) {
	GlobalDIDController.GetChallenge(c)
}

func DIDLogin(c *gin.Context) {
	GlobalDIDController.DIDLogin(c)
}

func IssueCredential(c *gin.Context) {
	GlobalDIDController.IssueCredential(c)
}

func GetCredentials(c *gin.Context) {
	GlobalDIDController.GetCredentials(c)
}

func VerifyPresentation(c *gin.Context) {
	GlobalDIDController.VerifyPresentation(c)
}

func RevokeCredential(c *gin.Context) {
	GlobalDIDController.RevokeCredential(c)
}

func DIDRegister(c *gin.Context) {
	GlobalDIDController.DIDRegister(c)
}

func GetDIDByUser(c *gin.Context) {
	GlobalDIDController.GetDIDByUser(c)
}
