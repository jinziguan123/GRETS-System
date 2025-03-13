package controller

import (
	"encoding/json"
	"grets_server/pkg/blockchain"
	"grets_server/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	RealtyID      string `json:"realtyId" binding:"required"`
	TransactionID string `json:"transactionId" binding:"required"`
	Seller        string `json:"seller" binding:"required"`
	Buyer         string `json:"buyer" binding:"required"`
	Content       string `json:"content" binding:"required"`
	ContractType  string `json:"contractType" binding:"required"`
	ValidityDate  string `json:"validityDate" binding:"required"`
}

// CreateContract 创建合同
func CreateContract(c *gin.Context) {
	// 解析请求参数
	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 获取当前用户
	userId := c.GetString("userId")
	if userId == "" {
		utils.ResponseUnauthorized(c, "未获取到用户信息")
		return
	}

	// 调用链码创建合同
	_, err := blockchain.DefaultFabricClient.Invoke("CreateContract",
		req.RealtyID,
		req.TransactionID,
		req.Seller,
		req.Buyer,
		req.Content,
		req.ContractType,
		req.ValidityDate,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "创建合同失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}

// QueryContractListRequest 查询合同列表请求
type QueryContractListRequest struct {
	RealtyID      string `form:"realtyId"`
	TransactionID string `form:"transactionId"`
	Seller        string `form:"seller"`
	Buyer         string `form:"buyer"`
	ContractType  string `form:"contractType"`
	Status        string `form:"status"`
	PageSize      int    `form:"pageSize,default=10"`
	PageNumber    int    `form:"pageNumber,default=1"`
}

// QueryContractList 查询合同列表
func QueryContractList(c *gin.Context) {
	// 解析查询参数
	var req QueryContractListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的查询参数")
		return
	}

	// 准备查询条件的JSON字符串
	condition := map[string]interface{}{}
	if req.RealtyID != "" {
		condition["realtyId"] = req.RealtyID
	}
	if req.TransactionID != "" {
		condition["transactionId"] = req.TransactionID
	}
	if req.Seller != "" {
		condition["seller"] = req.Seller
	}
	if req.Buyer != "" {
		condition["buyer"] = req.Buyer
	}
	if req.ContractType != "" {
		condition["contractType"] = req.ContractType
	}
	if req.Status != "" {
		condition["status"] = req.Status
	}

	conditionJson, err := json.Marshal(condition)
	if err != nil {
		utils.ResponseInternalServerError(c, "序列化查询条件失败")
		return
	}

	// 调用链码查询合同列表
	result, err := blockchain.DefaultFabricClient.Query("QueryContractByCondition",
		string(conditionJson),
		strconv.Itoa(req.PageSize),
		strconv.Itoa(req.PageNumber))
	if err != nil {
		utils.ResponseInternalServerError(c, "查询合同列表失败")
		return
	}

	// 解析响应数据
	var resp struct {
		Total int         `json:"total"`
		List  interface{} `json:"list"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		utils.ResponseInternalServerError(c, "解析响应数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, resp)
}

// GetContractByID 根据ID获取合同信息
func GetContractByID(c *gin.Context) {
	// 获取合同ID
	contractId := c.Param("id")
	if contractId == "" {
		utils.ResponseBadRequest(c, "无效的合同ID")
		return
	}

	// 调用链码查询合同
	result, err := blockchain.DefaultFabricClient.Query("QueryContract", contractId)
	if err != nil {
		utils.ResponseInternalServerError(c, "查询合同失败")
		return
	}

	// 检查是否找到合同
	if len(result) == 0 {
		utils.ResponseNotFound(c, "合同不存在")
		return
	}

	// 解析合同数据
	var contract interface{}
	if err := json.Unmarshal(result, &contract); err != nil {
		utils.ResponseInternalServerError(c, "解析合同数据失败")
		return
	}

	// 返回结果
	utils.ResponseWithData(c, contract)
}

// SignContractRequest 签署合同请求
type SignContractRequest struct {
	SignatureInfo string `json:"signatureInfo" binding:"required"`
	SignatureHash string `json:"signatureHash" binding:"required"`
	SignatureDate string `json:"signatureDate" binding:"required"`
}

// SignContract 签署合同
func SignContract(c *gin.Context) {
	// 获取合同ID
	contractId := c.Param("id")
	if contractId == "" {
		utils.ResponseBadRequest(c, "无效的合同ID")
		return
	}

	// 解析请求参数
	var req SignContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseBadRequest(c, "无效的请求参数")
		return
	}

	// 获取当前用户
	userId := c.GetString("userId")
	if userId == "" {
		utils.ResponseUnauthorized(c, "未获取到用户信息")
		return
	}

	// 调用链码签署合同
	_, err := blockchain.DefaultFabricClient.Invoke("SignContract",
		contractId,
		req.SignatureInfo,
		req.SignatureHash,
		req.SignatureDate,
		userId)
	if err != nil {
		utils.ResponseInternalServerError(c, "签署合同失败: "+err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c)
}
