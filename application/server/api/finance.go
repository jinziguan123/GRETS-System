package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grets/server/sdk"
)

// FinanceHandler 金融API处理器
type FinanceHandler struct {
	FabricClient *sdk.FabricClient
}

// NewFinanceHandler 创建新的金融API处理器
func NewFinanceHandler(client *sdk.FabricClient) *FinanceHandler {
	return &FinanceHandler{
		FabricClient: client,
	}
}

// RegisterRoutes 注册路由
func (h *FinanceHandler) RegisterRoutes(router *gin.RouterGroup) {
	financeGroup := router.Group("/finance")
	{
		financeGroup.GET("/loan/list", h.ListLoans)
		financeGroup.GET("/loan/:id", h.GetLoan)
		financeGroup.POST("/loan/apply", h.ApplyLoan)
		financeGroup.POST("/loan/approve", h.ApproveLoan)
		financeGroup.POST("/loan/reject", h.RejectLoan)
		financeGroup.POST("/loan/disburse", h.DisburseLoan)
		financeGroup.POST("/loan/repay", h.RepayLoan)
		financeGroup.GET("/loan/user/:userId", h.GetLoansByUser)
		financeGroup.GET("/loan/property/:propertyId", h.GetLoansByProperty)
	}
}

// ListLoans 获取贷款列表
func (h *FinanceHandler) ListLoans(c *gin.Context) {
	// 获取查询参数
	status := c.Query("status")

	var result []byte
	var err error

	// 根据参数决定调用的链码函数
	if status != "" {
		// 查询特定状态的贷款
		result, err = h.FabricClient.Query("financechannel", "financecc", "GetLoansByStatus", status)
	} else {
		// 查询所有贷款
		result, err = h.FabricClient.Query("financechannel", "financecc", "GetAllLoans")
	}

	if err != nil {
		log.Printf("查询贷款列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询贷款列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var loans []map[string]interface{}
	if err := json.Unmarshal(result, &loans); err != nil {
		log.Printf("解析贷款列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析贷款列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   loans,
	})
}

// GetLoan 获取贷款详情
func (h *FinanceHandler) GetLoan(c *gin.Context) {
	loanID := c.Param("id")
	if loanID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "贷款ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("financechannel", "financecc", "GetLoan", loanID)
	if err != nil {
		log.Printf("查询贷款详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询贷款详情失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var loan map[string]interface{}
	if err := json.Unmarshal(result, &loan); err != nil {
		log.Printf("解析贷款详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析贷款详情失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   loan,
	})
}

// ApplyLoan 申请贷款
func (h *FinanceHandler) ApplyLoan(c *gin.Context) {
	var requestBody struct {
		LoanID          string `json:"loanId" binding:"required"`
		BorrowerID      string `json:"borrowerId" binding:"required"`
		PropertyID      string `json:"propertyId" binding:"required"`
		Amount          string `json:"amount" binding:"required"`
		Term            string `json:"term" binding:"required"` // 贷款期限（月）
		InterestRate    string `json:"interestRate" binding:"required"`
		Purpose         string `json:"purpose" binding:"required"`
		CollateralType  string `json:"collateralType"`
		CollateralValue string `json:"collateralValue"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码申请贷款
	result, err := h.FabricClient.Invoke(
		"financechannel",
		"financecc",
		"ApplyLoan",
		requestBody.LoanID,
		requestBody.BorrowerID,
		requestBody.PropertyID,
		requestBody.Amount,
		requestBody.Term,
		requestBody.InterestRate,
		requestBody.Purpose,
		requestBody.CollateralType,
		requestBody.CollateralValue,
	)

	if err != nil {
		log.Printf("申请贷款失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "申请贷款失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "贷款申请成功",
		"data":    string(result),
	})
}

// ApproveLoan 批准贷款
func (h *FinanceHandler) ApproveLoan(c *gin.Context) {
	var requestBody struct {
		LoanID     string `json:"loanId" binding:"required"`
		ApproverID string `json:"approverId" binding:"required"`
		Remarks    string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码批准贷款
	result, err := h.FabricClient.Invoke(
		"financechannel",
		"financecc",
		"ApproveLoan",
		requestBody.LoanID,
		requestBody.ApproverID,
		requestBody.Remarks,
	)

	if err != nil {
		log.Printf("批准贷款失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "批准贷款失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "贷款已批准",
		"data":    string(result),
	})
}

// RejectLoan 拒绝贷款
func (h *FinanceHandler) RejectLoan(c *gin.Context) {
	var requestBody struct {
		LoanID     string `json:"loanId" binding:"required"`
		ApproverID string `json:"approverId" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码拒绝贷款
	result, err := h.FabricClient.Invoke(
		"financechannel",
		"financecc",
		"RejectLoan",
		requestBody.LoanID,
		requestBody.ApproverID,
		requestBody.Reason,
	)

	if err != nil {
		log.Printf("拒绝贷款失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "拒绝贷款失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "贷款已拒绝",
		"data":    string(result),
	})
}

// DisburseLoan 发放贷款
func (h *FinanceHandler) DisburseLoan(c *gin.Context) {
	var requestBody struct {
		LoanID  string `json:"loanId" binding:"required"`
		BankID  string `json:"bankId" binding:"required"`
		Amount  string `json:"amount" binding:"required"`
		Remarks string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码发放贷款
	result, err := h.FabricClient.Invoke(
		"financechannel",
		"financecc",
		"DisburseLoan",
		requestBody.LoanID,
		requestBody.BankID,
		requestBody.Amount,
		requestBody.Remarks,
	)

	if err != nil {
		log.Printf("发放贷款失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "发放贷款失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "贷款已发放",
		"data":    string(result),
	})
}

// RepayLoan 还款
func (h *FinanceHandler) RepayLoan(c *gin.Context) {
	var requestBody struct {
		LoanID      string `json:"loanId" binding:"required"`
		BorrowerID  string `json:"borrowerId" binding:"required"`
		Amount      string `json:"amount" binding:"required"`
		PaymentType string `json:"paymentType" binding:"required"` // 本金、利息或全部
		Remarks     string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 调用链码还款
	result, err := h.FabricClient.Invoke(
		"financechannel",
		"financecc",
		"RepayLoan",
		requestBody.LoanID,
		requestBody.BorrowerID,
		requestBody.Amount,
		requestBody.PaymentType,
		requestBody.Remarks,
	)

	if err != nil {
		log.Printf("还款失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "还款失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "还款成功",
		"data":    string(result),
	})
}

// GetLoansByUser 获取用户的贷款记录
func (h *FinanceHandler) GetLoansByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "用户ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("financechannel", "financecc", "GetLoansByBorrower", userID)
	if err != nil {
		log.Printf("查询用户贷款记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询用户贷款记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var loans []map[string]interface{}
	if err := json.Unmarshal(result, &loans); err != nil {
		log.Printf("解析用户贷款记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析用户贷款记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   loans,
	})
}

// GetLoansByProperty 获取特定房产的贷款记录
func (h *FinanceHandler) GetLoansByProperty(c *gin.Context) {
	propertyID := c.Param("propertyId")
	if propertyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "房产ID不能为空",
		})
		return
	}

	result, err := h.FabricClient.Query("financechannel", "financecc", "GetLoansByProperty", propertyID)
	if err != nil {
		log.Printf("查询房产贷款记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询房产贷款记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析结果
	var loans []map[string]interface{}
	if err := json.Unmarshal(result, &loans); err != nil {
		log.Printf("解析房产贷款记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "解析房产贷款记录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   loans,
	})
}
