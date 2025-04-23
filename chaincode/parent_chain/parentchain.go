package main

import (
	"encoding/json"
	"fmt"
	"log"
	"parent_chain_chaincode/constances"
	"parent_chain_chaincode/models"
	"parent_chain_chaincode/tools"
	"time"

	"maps"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 抵押信息结构
type Mortgage struct {
	MortgageID      string    `json:"mortgageId"`      // 抵押ID
	RealEstateID    string    `json:"realEstateId"`    // 房产ID
	BankID          string    `json:"bankId"`          // 银行ID
	BorrowerID      string    `json:"borrowerId"`      // 借款人ID
	LoanAmount      float64   `json:"loanAmount"`      // 贷款金额
	InterestRate    float64   `json:"interestRate"`    // 利率
	Term            int       `json:"term"`            // 期限(月)
	StartDate       time.Time `json:"startDate"`       // 开始日期
	EndDate         time.Time `json:"endDate"`         // 结束日期
	Status          string    `json:"status"`          // 状态
	ApprovedBy      string    `json:"approvedBy"`      // 批准人
	ApprovedAt      time.Time `json:"approvedAt"`      // 批准时间
	LastUpdateTime  time.Time `json:"lastUpdateTime"`  // 最后更新时间
	PaymentPlan     string    `json:"paymentPlan"`     // 还款计划
	CollateralValue float64   `json:"collateralValue"` // 抵押物估值
}

// 税费信息结构
type Tax struct {
	TaxID         string    `json:"taxId"`         // 税费ID
	TransactionID string    `json:"transactionId"` // 关联交易ID
	TaxType       string    `json:"taxType"`       // 税费类型
	TaxRate       float64   `json:"taxRate"`       // 税率
	TaxAmount     float64   `json:"taxAmount"`     // 税额
	Status        string    `json:"status"`        // 状态（已缴/未缴）
	DueDate       time.Time `json:"dueDate"`       // 截止日期
	PaidAt        time.Time `json:"paidAt"`        // 缴纳时间
	PaidBy        string    `json:"paidBy"`        // 缴纳人
	ReceiptID     string    `json:"receiptId"`     // 收据ID
}

// 审计记录结构
type AuditRecord struct {
	AuditID         string    `json:"auditId"`         // 审计ID
	TargetType      string    `json:"targetType"`      // 目标类型(房产/交易/抵押)
	TargetID        string    `json:"targetId"`        // 目标ID
	AuditorID       string    `json:"auditorId"`       // 审计员ID
	AuditorOrgID    string    `json:"auditorOrgId"`    // 审计员组织ID
	Result          string    `json:"result"`          // 审计结果
	Comments        string    `json:"comments"`        // 审计意见
	AuditedAt       time.Time `json:"auditedAt"`       // 审计时间
	Violations      []string  `json:"violations"`      // 违规项
	Recommendations []string  `json:"recommendations"` // 建议
}

// 查询结果结构
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// 获取调用者身份MSP ID
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("[getClientIdentityMSPID] 获取客户端MSP ID失败: %v", err)
	}

	return clientMSPID, nil

}

// 创建复合键
func (s *SmartContract) createCompositeKey(ctx contractapi.TransactionContextInterface, objectType string,
	attributes ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("[createCompositeKey] 创建复合键失败: %v", err)
	}
	return key, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 用户相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetUserByCitizenIDAndOrganization 根据身份证号和组织获取用户信息
func (s *SmartContract) GetUserByCitizenIDAndOrganization(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	organization string,
) (*models.UserPublic, error) {
	// 检查必填参数
	if len(citizenIDHash) == 0 {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 身份证号不能为空")
	}
	if len(organization) == 0 {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 组织不能为空")
	}

	// 生成复合键：身份证号-组织
	key, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 创建复合键失败: %v", err)
	}

	// 通过复合键查找用户ID
	userBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 查询用户ID失败: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 用户不存在")
	}

	// 解析用户数据
	var user models.UserPublic
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByCitizenIDAndOrganization] 解析用户数据失败: %v", err)
	}

	return &user, nil
}

// GetBalanceByCitizenIDHashAndOrganization 根据身份证号和组织获取用户余额
func (s *SmartContract) GetBalanceByCitizenIDHashAndOrganization(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	organization string,
) (float64, error) {
	key, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return 0, fmt.Errorf("[GetBalanceByCitizenIDHashAndOrganization] 创建复合键失败: %v", err)
	}

	userBytes, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, key)
	if err != nil {
		return 0, fmt.Errorf("[GetBalanceByCitizenIDHashAndOrganization] 查询用户数据失败: %v", err)
	}
	if userBytes == nil {
		return 0, fmt.Errorf("[GetBalanceByCitizenIDHashAndOrganization] 用户不存在")
	}

	var user models.UserPrivate
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return 0, fmt.Errorf("[GetBalanceByCitizenIDHashAndOrganization] 解析用户数据失败: %v", err)
	}

	return user.Balance, nil
}

// UpdateUser 更新用户信息
func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	organization string,
	name string,
	phone string,
	email string,
	passwordHash string,
	status string,
) error {
	key, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{citizenIDHash}...)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 创建复合键失败: %v", err)
	}

	// 获取现有用户数据
	userPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 查询用户公开信息失败: %v", err)
	}
	if userPublicBytes == nil {
		return fmt.Errorf("[UpdateUser] 用户不存在")
	}

	userPrivateBytes, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, key)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 查询用户私密信息失败: %v", err)
	}
	if userPrivateBytes == nil {
		return fmt.Errorf("[UpdateUser] 用户不存在")
	}

	// 解析用户数据
	var userPublic models.UserPublic
	err = json.Unmarshal(userPublicBytes, &userPublic)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 解析用户公开信息失败: %v", err)
	}

	var userPrivate models.UserPrivate
	err = json.Unmarshal(userPrivateBytes, &userPrivate)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 解析用户私密信息失败: %v", err)
	}

	// 更新用户数据
	if len(name) > 0 {
		userPublic.Name = name
	}
	if len(phone) > 0 {
		userPrivate.Phone = phone
	}
	if len(email) > 0 {
		userPrivate.Email = email
	}
	if len(passwordHash) > 0 {
		userPrivate.PasswordHash = passwordHash
	}
	if len(status) > 0 {
		userPublic.Status = status
	}
	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateUser] 获取交易时间戳失败: %v", err)
	}
	userPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化并保存更新后的用户数据
	updatedUserPublicBytes, err := json.Marshal(userPublic)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 序列化用户数据失败: %v", err)
	}

	err = ctx.GetStub().PutState(key, updatedUserPublicBytes)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 保存用户数据失败: %v", err)
	}

	updatedUserPrivateBytes, err := json.Marshal(userPrivate)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 序列化用户数据失败: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(constances.UserDataCollection, key, updatedUserPrivateBytes)
	if err != nil {
		return fmt.Errorf("[UpdateUser] 保存用户数据失败: %v", err)
	}

	return nil
}

// Register 注册用户
func (s *SmartContract) Register(ctx contractapi.TransactionContextInterface,
	citizenIDHash string,
	citizenID string,
	name string,
	phone string,
	email string,
	passwordHash string,
	organization string,
	role string,
	status string,
	balance float64,
) error {
	// 检查用户是否已存在
	key, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return fmt.Errorf("[Register] 创建复合键失败: %v", err)
	}
	existUser, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[Register] 查询用户失败: %v", err)
	}
	if existUser != nil {
		return fmt.Errorf("[Register] 用户已存在")
	}

	userKey, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{citizenIDHash, organization}...)
	if err != nil {
		return fmt.Errorf("[Register] 创建复合键失败: %v", err)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[Register] 获取交易时间戳失败: %v", err)
	}

	// 创建用户公开信息
	userPublic := models.UserPublic{
		DocType:        constances.DocTypeUser,
		CitizenID:      citizenID,
		Name:           name,
		Organization:   organization,
		Role:           role,
		Status:         status,
		CreateTime:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		LastUpdateTime: time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化用户
	userJSON, err := json.Marshal(userPublic)
	if err != nil {
		return fmt.Errorf("[Register] 序列化用户失败: %v", err)
	}

	// 保存用户
	err = ctx.GetStub().PutState(userKey, userJSON)
	if err != nil {
		return fmt.Errorf("[Register] 保存用户失败: %v", err)
	}

	// 创建用户私钥
	userPrivate := models.UserPrivate{
		DocType:      constances.DocTypeUser,
		CitizenID:    citizenID,
		PasswordHash: passwordHash,
		Balance:      balance,
		Phone:        phone,
		Email:        email,
	}

	// 序列化用户私钥
	userPrivateJSON, err := json.Marshal(userPrivate)
	if err != nil {
		return fmt.Errorf("[Register] 序列化用户私钥失败: %v", err)
	}

	// 保存用户私钥
	err = ctx.GetStub().PutPrivateData(constances.UserDataCollection, userKey, userPrivateJSON)
	if err != nil {
		return fmt.Errorf("[Register] 保存用户私钥失败: %v", err)
	}
	return nil
}

// ListUsersByOrganization 查询特定组织的用户
func (s *SmartContract) ListUsersByOrganization(ctx contractapi.TransactionContextInterface,
	organization string,
) ([]*models.UserPublic, error) {
	// 获取所有用户
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(constances.DocTypeUser, []string{organization})
	if err != nil {
		return nil, fmt.Errorf("[ListUsersByOrganization] 创建迭代器失败: %v", err)
	}
	defer resultsIterator.Close()

	var users []*models.UserPublic
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 获取下一个用户失败: %v", err)
		}

		key := queryResponse.Key
		userPublic, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 查询用户公开信息失败: %v", err)
		}

		userPrivate, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, key)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 查询用户私钥失败: %v", err)
		}

		var userPublicMap, userPrivateMap map[string]interface{}
		err = json.Unmarshal(userPublic, &userPublicMap)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 解析用户公开信息失败: %v", err)
		}

		err = json.Unmarshal(userPrivate, &userPrivateMap)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 解析用户私钥失败: %v", err)
		}

		mergedMap := make(map[string]interface{})
		for k, v := range userPublicMap {
			mergedMap[k] = v
		}
		for k, v := range userPrivateMap {
			mergedMap[k] = v
		}

		mergedJSON, err := json.Marshal(mergedMap)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 序列化用户失败: %v", err)
		}

		var user models.UserPublic
		err = json.Unmarshal(mergedJSON, &user)
		if err != nil {
			return nil, fmt.Errorf("[ListUsersByOrganization] 序列化用户失败: %v", err)
		}

		// 检查是否属于指定组织
		if user.Organization == organization {
			users = append(users, &user)
		}
	}

	return users, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 房产相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateRealty 创建房产信息（仅政府机构可调用）
func (s *SmartContract) CreateRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	realtyCert string,
	realtyType string,
	currentOwnerCitizenIDHash string,
	currentOwnerOrganization string,
	previousOwnersCitizenIDHashListJSON string,
) error {
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("[CreateRealty] 获取客户端ID失败: %v", err)
	}

	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.GovernmentMSP {
		return fmt.Errorf("[CreateRealty] 只有政府机构可以创建房产信息")
	}

	// 创建房产信息复合键
	key, err := s.createCompositeKey(ctx, constances.DocTypeRealEstate, []string{realtyCertHash}...)
	if err != nil {
		return err
	}

	// 检查房产是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 查询房产信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("[CreateRealty] 房产ID %s 已存在", realtyCertHash)
	}

	// 解析JSON字符串为字符串数组
	var previousOwnersCitizenIDHashList []string
	if err := json.Unmarshal([]byte(previousOwnersCitizenIDHashListJSON), &previousOwnersCitizenIDHashList); err != nil {
		return fmt.Errorf("[CreateRealty] 解析历史所有者列表失败: %v", err)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreateRealty] 获取交易时间戳失败: %v", err)
	}

	// 创建公开房产信息
	realEstate := models.Realty{
		DocType:                         constances.DocTypeRealEstate,
		RealtyCertHash:                  realtyCertHash,
		RealtyCert:                      realtyCert,
		RealtyType:                      realtyType,
		CreateTime:                      time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		Status:                          constances.RealtyStatusNormal,
		LastUpdateTime:                  time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		CurrentOwnerCitizenIDHash:       currentOwnerCitizenIDHash,
		CurrentOwnerOrganization:        currentOwnerOrganization,
		PreviousOwnersCitizenIDHashList: previousOwnersCitizenIDHashList,
	}

	// 序列化并保存
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 序列化公开房产信息失败: %v", err)
	}

	ctx.GetStub().PutState(key, realEstateJSON)

	// 创建房产私钥
	realEstatePrivate := models.RealtyPrivate{
		DocType:                         constances.DocTypeRealEstate,
		RealtyCertHash:                  realtyCertHash,
		RealtyCert:                      realtyCert,
		CurrentOwnerCitizenIDHash:       currentOwnerCitizenIDHash,
		CurrentOwnerOrganization:        currentOwnerOrganization,
		PreviousOwnersCitizenIDHashList: previousOwnersCitizenIDHashList,
	}

	// 序列化并保存
	realEstatePrivateJSON, err := json.Marshal(realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 序列化房产私钥失败: %v", err)
	}

	err = ctx.GetStub().PutPrivateData(constances.RealEstatePrivateCollection, key, realEstatePrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 保存房产私钥失败: %v", err)
	}

	// 创建房产登记记录
	key, err = s.createCompositeKey(ctx, constances.DocTypeRealEstate, []string{realtyCertHash, "createRealty"}...)
	if err != nil {
		return err
	}
	type createRealtyRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := createRealtyRecord{
		ClientID: clientID,
		Action:   "createRealty",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CreateRealty] 序列化房产登记记录失败: %v", err)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) QueryRealtyByOrganizationAndCitizenIDHash(ctx contractapi.TransactionContextInterface,
	organization string,
	citizenIDHash string,
) ([]*models.Realty, error) {
	// 构建查询语句
	queryString := fmt.Sprintf(`{
		"selector": {
			"docType": "%s",
			"currentOwnerOrganization": "%s",
			"currentOwnerCitizenIDHash": "%s"
		}
	}`, constances.DocTypeRealEstate, organization, citizenIDHash)

	realties, err := tools.SelectByQueryString[models.Realty](ctx, queryString)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealtyByOrganizationAndCitizenIDHash] 查询房产信息失败: %v", err)
	}

	realtyList := []*models.Realty{}
	realtyList = append(realtyList, realties...)

	return realtyList, nil
}

// QueryRealty 查询房产信息
func (s *SmartContract) QueryRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
) (*models.Realty, error) {

	key, err := s.createCompositeKey(ctx, constances.DocTypeRealEstate, []string{realtyCertHash}...)
	if err != nil {
		return nil, err
	}

	realEstatePublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 查询房产信息失败: %v", err)
	}
	if realEstatePublicBytes == nil {
		return nil, fmt.Errorf("[QueryRealty] 房产ID %s 不存在", realtyCertHash)
	}

	var realEstatePublic models.RealtyPublic
	err = json.Unmarshal(realEstatePublicBytes, &realEstatePublic)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 解析房产信息失败: %v", err)
	}

	var realEstatePrivate models.RealtyPrivate
	realEstatePrivateBytes, err := ctx.GetStub().GetPrivateData(constances.RealEstatePrivateCollection, key)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 查询房产私钥失败: %v", err)
	}
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivate)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 解析房产私钥失败: %v", err)
	}

	var realEstatePublicMap, realEstatePrivateMap map[string]interface{}
	err = json.Unmarshal(realEstatePublicBytes, &realEstatePublicMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 解析房产信息失败: %v", err)
	}
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivateMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 解析房产私钥失败: %v", err)
	}

	realEstateMap := make(map[string]interface{})
	maps.Copy(realEstateMap, realEstatePublicMap)
	maps.Copy(realEstateMap, realEstatePrivateMap)
	realEstateJSON, err := json.Marshal(realEstateMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 序列化房产信息失败: %v", err)
	}

	var realEstate models.Realty
	err = json.Unmarshal(realEstateJSON, &realEstate)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealty] 解析房产信息失败: %v", err)
	}

	return &realEstate, nil
}

// QueryRealtyList 查询全量房产列表
func (s *SmartContract) QueryRealtyList(ctx contractapi.TransactionContextInterface,
	pageSize int32,
	bookmark string,
) ([]*models.RealtyPublic, error) {

	iter, _, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
		constances.DocTypeRealEstate,
		[]string{},
		pageSize,
		bookmark,
	)
	if err != nil {
		return nil, fmt.Errorf("[QueryRealtyList] 查询房产列表失败: %v", err)
	}
	defer iter.Close()

	realtyList := []*models.RealtyPublic{}
	for iter.HasNext() {
		realty, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("[QueryRealtyList] 查询房产列表失败: %v", err)
		}
		var realtyPublic models.RealtyPublic
		err = json.Unmarshal(realty.Value, &realtyPublic)
		if err != nil {
			return nil, fmt.Errorf("[QueryRealtyList] 解析房产信息失败: %v", err)
		}
		realtyList = append(realtyList, &realtyPublic)
	}
	return realtyList, nil
}

// UpdateRealty 更新房产信息（仅政府机构、投资者可调用）
func (s *SmartContract) UpdateRealty(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	realtyType string,
	status string,
	currentOwnerCitizenIDHash string,
	currentOwnerOrganization string,
	previousOwnersCitizenIDHashListJSON string,
) error {
	// 检查调用者身份
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 获取客户端ID失败: %v", err)
	}

	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.GovernmentMSP && clientMSPID != constances.InvestorMSP {
		return fmt.Errorf("[UpdateRealty] 只有政府机构、投资者可以更新房产信息")
	}

	// 查询现有房产信息
	realEstatePublic, err := s.QueryRealty(ctx, realtyCertHash)
	if err != nil {
		return err
	}

	if realEstatePublic.Status == constances.RealtyStatusFrozen {
		return fmt.Errorf("[UpdateRealty] 房产已被冻结，无法更新")
	}

	// 获取复合键
	key, err := s.createCompositeKey(ctx, constances.DocTypeRealEstate, []string{realEstatePublic.RealtyCertHash}...)
	if err != nil {
		return err
	}

	// 查询房产私钥
	realEstatePrivateBytes, err := ctx.GetStub().GetPrivateData(constances.RealEstatePrivateCollection, key)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 查询房产私钥失败: %v", err)
	}
	if realEstatePrivateBytes == nil {
		return fmt.Errorf("[UpdateRealty] 房产私钥不存在")
	}

	var realEstatePrivate models.RealtyPrivate
	err = json.Unmarshal(realEstatePrivateBytes, &realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 解析房产私钥失败: %v", err)
	}

	var modifyFields []string
	// 更新信息
	if realtyType != "" && realtyType != realEstatePublic.RealtyType {
		realEstatePublic.RealtyType = realtyType
		modifyFields = append(modifyFields, "realtyType")
	}
	if status != "" && status != realEstatePublic.Status {
		realEstatePublic.Status = status
		modifyFields = append(modifyFields, "status")
	}
	if currentOwnerCitizenIDHash != "" && currentOwnerCitizenIDHash != realEstatePrivate.CurrentOwnerCitizenIDHash {
		realEstatePrivate.CurrentOwnerCitizenIDHash = currentOwnerCitizenIDHash
		modifyFields = append(modifyFields, "currentOwnerCitizenIDHash")
	}
	if currentOwnerOrganization != "" && currentOwnerOrganization != realEstatePrivate.CurrentOwnerOrganization {
		realEstatePrivate.CurrentOwnerOrganization = currentOwnerOrganization
	}
	// 解析JSON字符串为字符串数组
	var previousOwnersCitizenIDHashList []string
	if err := json.Unmarshal([]byte(previousOwnersCitizenIDHashListJSON), &previousOwnersCitizenIDHashList); err != nil {
		return fmt.Errorf("[UpdateRealty] 解析历史所有者列表失败: %v", err)
	}
	if len(previousOwnersCitizenIDHashList) > len(realEstatePrivate.PreviousOwnersCitizenIDHashList) {
		realEstatePrivate.PreviousOwnersCitizenIDHashList = previousOwnersCitizenIDHashList
		modifyFields = append(modifyFields, "previousOwnersCitizenIDHashList")
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 获取交易时间戳失败: %v", err)
	}
	realEstatePublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化并保存
	realEstatePublicJSON, err := json.Marshal(realEstatePublic)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 序列化房产信息失败: %v", err)
	}

	ctx.GetStub().PutState(key, realEstatePublicJSON)

	// 私有数据序列化并保存
	realEstatePrivateJSON, err := json.Marshal(realEstatePrivate)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 序列化房产私钥失败: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(constances.RealEstatePrivateCollection, key, realEstatePrivateJSON)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 保存房产私钥失败: %v", err)
	}

	// 创建房产登记记录
	key, err = s.createCompositeKey(ctx, constances.DocTypeRealEstate, []string{realEstatePublic.RealtyCertHash, "updateRealty"}...)
	if err != nil {
		return err
	}

	type updateRealtyRecord struct {
		ClientID     string    `json:"clientID"`
		Action       string    `json:"action"`
		Time         time.Time `json:"time"`
		ModifyFields []string  `json:"modifyFields"`
	}
	record := updateRealtyRecord{
		ClientID:     clientID,
		Action:       "updateRealty",
		Time:         time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		ModifyFields: modifyFields,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[UpdateRealty] 序列化房产登记记录失败: %v", err)
	}
	return ctx.GetStub().PutState(key, recordJSON)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 交易相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateTransaction 创建交易（投资者、政府可以调用）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface,
	realtyCertHash string,
	transactionUUID string,
	sellerCitizenIDHash string,
	sellerOrganization string,
	buyerCitizenIDHash string,
	buyerOrganization string,
	contractUUID string,
	paymentUUIDListJSON string,
	tax float64,
	price float64,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP {
		return fmt.Errorf("[CreateTransaction] 只有投资者、政府可以创建交易")
	}

	// 查询房产信息
	realEstate, err := s.QueryRealty(ctx, realtyCertHash)
	if err != nil {
		return err
	}

	// 检查房产状态
	if realEstate.Status != constances.RealtyStatusPendingSale {
		return fmt.Errorf("[CreateTransaction] 房产状态不允许交易: %s", realEstate.Status)
	}

	// 检查卖方是否为房产所有者
	if realEstate.CurrentOwnerCitizenIDHash != sellerCitizenIDHash {
		return fmt.Errorf("[CreateTransaction] 卖方不是房产所有者")
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 获取交易时间戳失败: %v", err)
	}

	// 创建公开交易信息
	transactionPublic := models.TransactionPublic{
		DocType:             constances.DocTypeTransaction,
		TransactionUUID:     transactionUUID,
		RealtyCertHash:      realtyCertHash,
		SellerCitizenIDHash: sellerCitizenIDHash,
		SellerOrganization:  sellerOrganization,
		BuyerCitizenIDHash:  buyerCitizenIDHash,
		BuyerOrganization:   buyerOrganization,
		Status:              constances.TxStatusPending,
		CreateTime:          time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		UpdateTime:          time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 解析JSON字符串为字符串数组
	var paymentUUIDList []string
	if err := json.Unmarshal([]byte(paymentUUIDListJSON), &paymentUUIDList); err != nil {
		return fmt.Errorf("[CreateTransaction] 解析支付UUID列表失败: %v", err)
	}

	// 创建私有交易信息
	transactionPrivate := models.TransactionPrivate{
		DocType:         constances.DocTypeTransaction,
		TransactionUUID: transactionUUID,
		Price:           price,
		Tax:             tax,
		PaymentUUIDList: paymentUUIDList,
		ContractUUID:    contractUUID,
	}

	// 创建交易复合键
	key, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	// 序列化并保存交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 序列化公开交易信息失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 保存公开交易信息失败: %v", err)
	}

	transactionPrivateJSON, err := json.Marshal(transactionPrivate)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 序列化私有交易信息失败: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(constances.TransactionPrivateCollection, key, transactionPrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 保存私有交易信息失败: %v", err)
	}

	// 创建交易登记记录
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}

	key, err = s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID, "createTransaction"}...)
	if err != nil {
		return err
	}
	type createTransactionRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := createTransactionRecord{
		ClientID: clientID,
		Action:   "createTransaction",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 序列化交易登记记录失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("[CreateTransaction] 保存交易登记记录失败: %v", err)
	}

	return nil
}

// QueryTransaction 查询交易（投资者、政府可以调用）
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
) (*models.Transaction, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 获取客户端ID失败: %v", err)
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP {
		return nil, fmt.Errorf("[QueryTransaction] 只有投资者、政府可以查询交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 创建复合键失败: %v", err)
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 查询交易信息失败: %v", err)
	}

	if transactionPublicBytes == nil {
		return nil, fmt.Errorf("[QueryTransaction] 交易不存在: %s", transactionUUID)
	}

	transactionPrivateBytes, err := ctx.GetStub().GetPrivateData(constances.TransactionPrivateCollection, key)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 查询交易私钥失败: %v", err)
	}
	if transactionPrivateBytes == nil {
		return nil, fmt.Errorf("[QueryTransaction] 交易私钥不存在")
	}

	var transactionPublicMap, transactionPrivateMap map[string]interface{}
	err = json.Unmarshal(transactionPublicBytes, &transactionPublicMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 解析交易公开信息失败: %v", err)
	}

	err = json.Unmarshal(transactionPrivateBytes, &transactionPrivateMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 解析交易私钥失败: %v", err)
	}

	mergedMap := make(map[string]interface{})
	for k, v := range transactionPublicMap {
		mergedMap[k] = v
	}
	for k, v := range transactionPrivateMap {
		mergedMap[k] = v
	}

	mergedJSON, err := json.Marshal(mergedMap)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 序列化交易失败: %v", err)
	}

	var transaction models.Transaction
	err = json.Unmarshal(mergedJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransaction] 序列化交易失败: %v", err)
	}

	return &transaction, nil
}

func (s *SmartContract) QueryTransactionList(ctx contractapi.TransactionContextInterface,
	pageSize int32,
	bookmark string,
) ([]*models.TransactionPublic, error) {

	iter, _, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
		constances.DocTypeTransaction,
		[]string{},
		pageSize,
		bookmark,
	)
	if err != nil {
		return nil, fmt.Errorf("[QueryTransactionList] 查询交易列表失败: %v", err)
	}
	defer iter.Close()

	var transactionList []*models.TransactionPublic
	for iter.HasNext() {
		transaction, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("[QueryTransactionList] 查询交易列表失败: %v", err)
		}
		var transactionPublic models.TransactionPublic
		err = json.Unmarshal(transaction.Value, &transactionPublic)
		if err != nil {
			return nil, fmt.Errorf("[QueryTransactionList] 解析交易信息失败: %v", err)
		}
		transactionList = append(transactionList, &transactionPublic)
	}
	return transactionList, nil
}

// CheckTransaction 检查交易（投资者、政府可以调用）
func (s *SmartContract) CheckTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	status string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP {
		return fmt.Errorf("[CheckTransaction] 只有投资者、政府可以检查交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 查询交易信息失败: %v", err)
	}

	if transactionPublicBytes == nil {
		return fmt.Errorf("[CheckTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic models.TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 解析交易信息失败: %v", err)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 获取交易时间戳失败: %v", err)
	}

	// 检查交易状态
	if transactionPublic.Status != constances.TxStatusPending {
		return fmt.Errorf("[CheckTransaction] 交易状态不允许检查: %s", transactionPublic.Status)
	}

	// 更新交易状态
	transactionPublic.Status = status
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()
	// 序列化交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易信息
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 保存交易信息失败: %v", err)
	}

	// 创建交易审核记录
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	key, err = s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID, "checkTransaction"}...)
	if err != nil {
		return err
	}
	type checkTransactionRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := checkTransactionRecord{
		ClientID: clientID,
		Action:   "checkTransaction",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 序列化交易登记记录失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("[CheckTransaction] 保存交易登记记录失败: %v", err)
	}

	return nil
}

// UpdateTransaction 更新交易（投资者、政府可以调用）
func (s *SmartContract) UpdateTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	status string,
) error {

	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP {
		return fmt.Errorf("[UpdateTransaction] 只有投资者、政府可以更新交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 查询交易信息失败: %v", err)
	}

	if transactionPublicBytes == nil {
		return fmt.Errorf("[UpdateTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic models.TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 解析交易信息失败: %v", err)
	}

	if transactionPublic.Status != constances.TxStatusPending && transactionPublic.Status != constances.TxStatusInProgress {
		return fmt.Errorf("[UpdateTransaction] 交易状态不允许更新: %s", transactionPublic.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 获取交易时间戳失败: %v", err)
	}

	// 更新交易状态
	transactionPublic.Status = status
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易信息
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[UpdateTransaction] 保存交易信息失败: %v", err)
	}

	return nil
}

// CompleteTransaction 完成交易（投资者、政府可以调用）
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP {
		return fmt.Errorf("[CompleteTransaction] 只有投资者、政府可以完成交易")
	}

	// 查询交易信息
	key, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return err
	}

	transactionPublicBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 查询交易信息失败: %v", err)
	}
	if transactionPublicBytes == nil {
		return fmt.Errorf("[CompleteTransaction] 交易不存在: %s", transactionUUID)
	}

	var transactionPublic models.TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 解析交易信息失败: %v", err)
	}

	if transactionPublic.Status != constances.TxStatusInProgress {
		return fmt.Errorf("[CompleteTransaction] 交易状态不允许完成: %s", transactionPublic.Status)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 获取交易时间戳失败: %v", err)
	}

	// 更新交易状态
	transactionPublic.Status = constances.TxStatusCompleted
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化交易信息
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易信息
	err = ctx.GetStub().PutState(key, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 保存交易信息失败: %v", err)
	}

	// 更新房产信息
	realtyIDHash := transactionPublic.RealtyCertHash
	realEstate, err := s.QueryRealty(ctx, realtyIDHash)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 查询房产信息失败: %v", err)
	}

	// 更新房产信息
	previousOwnersCitizenIDHashList := append(realEstate.PreviousOwnersCitizenIDHashList, transactionPublic.SellerCitizenIDHash)
	previousOwnersCitizenIDHashListJSON, err := json.Marshal(previousOwnersCitizenIDHashList)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 序列化历史所有者列表失败: %v", err)
	}
	err = s.UpdateRealty(
		ctx,
		realtyIDHash,
		realEstate.RealtyType,
		constances.RealtyStatusNormal,
		transactionPublic.BuyerCitizenIDHash,
		transactionPublic.BuyerOrganization,
		string(previousOwnersCitizenIDHashListJSON),
	)
	if err != nil {
		return fmt.Errorf("更新房产信息失败: %v", err)
	}

	// 创建交易完成记录
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	key, err = s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID, "completeTransaction"}...)
	if err != nil {
		return err
	}
	type completeTransactionRecord struct {
		ClientID string    `json:"clientID"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}
	record := completeTransactionRecord{
		ClientID: clientID,
		Action:   "completeTransaction",
		Time:     time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 序列化交易完成记录失败: %v", err)
	}
	err = ctx.GetStub().PutState(key, recordJSON)
	if err != nil {
		return fmt.Errorf("[CompleteTransaction] 保存交易完成记录失败: %v", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 支付相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreatePayment 创建支付信息（仅银行和投资者可调用）
func (s *SmartContract) CreatePayment(ctx contractapi.TransactionContextInterface,
	paymentUUID string,
	amount float64,
	fromCitizenIDHash string,
	fromOrganization string,
	toCitizenIDHash string,
	toOrganization string,
	paymentType string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.BankMSP && clientMSPID != constances.InvestorMSP {
		return fmt.Errorf("[CreatePayment] 只有银行和投资者可以创建支付信息")
	}

	// 检查支付信息是否已存在
	paymentKey, err := s.createCompositeKey(ctx, constances.DocTypePayment, []string{paymentUUID}...)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 创建复合键失败: %v", err)
	}
	paymentBytes, err := ctx.GetStub().GetState(paymentKey)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询支付信息失败: %v", err)
	}
	if paymentBytes != nil {
		return fmt.Errorf("[CreatePayment] 支付信息已存在: %s", paymentUUID)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreatePayment] 获取交易时间戳失败: %v", err)
	}

	// 创建支付信息
	payment := models.Payment{
		DocType:               constances.DocTypePayment,
		PaymentUUID:           paymentUUID,
		Amount:                amount,
		PayerCitizenIDHash:    fromCitizenIDHash,
		ReceiverCitizenIDHash: toCitizenIDHash,
		PaymentType:           paymentType,
		CreateTime:            time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化支付信息失败: %v", err)
	}

	// 保存支付信息
	err = ctx.GetStub().PutState(paymentKey, paymentJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存支付信息失败: %v", err)
	}

	fromUserKey, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{fromCitizenIDHash, fromOrganization}...)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 创建复合键失败: %v", err)
	}

	// 查询fromCitizenIDHash的余额
	fromCitizenPrivateBytes, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, fromUserKey)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户余额失败: %v", err)
	}
	if fromCitizenPrivateBytes == nil {
		return fmt.Errorf("[CreatePayment] 来源用户不存在: %s", fromCitizenIDHash)
	}

	var fromCitizenPrivate models.UserPrivate
	err = json.Unmarshal(fromCitizenPrivateBytes, &fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户余额失败: %v", err)
	}

	if fromCitizenPrivate.Balance < amount {
		return fmt.Errorf("[CreatePayment] 余额不足: %s", fromCitizenIDHash)
	}

	// 更新fromCitizenIDHash的余额
	fromCitizenPrivate.Balance -= amount

	// 序列化fromCitizenIDHash的余额
	fromCitizenPrivateJSON, err := json.Marshal(fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户余额失败: %v", err)
	}

	// 保存fromCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(constances.UserDataCollection, fromUserKey, fromCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户余额失败: %v", err)
	}

	// 更新fromCitizenIDHash的
	fromCitizenPublicBytes, err := ctx.GetStub().GetState(fromUserKey)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户信息失败: %v", err)
	}

	var fromCitizenPublic models.UserPublic
	err = json.Unmarshal(fromCitizenPublicBytes, &fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户信息失败: %v", err)
	}

	fromCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化fromCitizenIDHash的
	fromCitizenPublicJSON, err := json.Marshal(fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户信息失败: %v", err)
	}

	// 保存fromCitizenIDHash的
	err = ctx.GetStub().PutState(fromUserKey, fromCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户信息失败: %v", err)
	}

	toUserKey, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{toCitizenIDHash, toOrganization}...)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 创建复合键失败: %v", err)
	}

	// 查询toCitizenIDHash的余额
	toCitizenBytes, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, toUserKey)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户余额失败: %v", err)
	}
	if toCitizenBytes == nil {
		return fmt.Errorf("[CreatePayment] 目标用户不存在: %s", toCitizenIDHash)
	}

	var toCitizenPrivate models.UserPrivate
	err = json.Unmarshal(toCitizenBytes, &toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的余额
	toCitizenPrivate.Balance += amount

	// 序列化toCitizenIDHash的余额
	toCitizenPrivateJSON, err := json.Marshal(toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户余额失败: %v", err)
	}

	// 保存toCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(constances.UserDataCollection, toUserKey, toCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的
	toCitizenPublicBytes, err := ctx.GetStub().GetState(toUserKey)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 查询用户信息失败: %v", err)
	}

	var toCitizenPublic models.UserPublic
	err = json.Unmarshal(toCitizenPublicBytes, &toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 解析用户信息失败: %v", err)
	}

	toCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化toCitizenIDHash的
	toCitizenPublicJSON, err := json.Marshal(toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 序列化用户信息失败: %v", err)
	}

	// 保存toCitizenIDHash的
	err = ctx.GetStub().PutState(toUserKey, toCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[CreatePayment] 保存用户信息失败: %v", err)
	}

	return nil
}

// QueryPayment 查询支付信息
func (s *SmartContract) QueryPayment(ctx contractapi.TransactionContextInterface,
	paymentUUID string,
) (*models.Payment, error) {
	paymentKey, err := s.createCompositeKey(ctx, constances.DocTypePayment, []string{paymentUUID}...)
	if err != nil {
		return nil, fmt.Errorf("[QueryPayment] 创建复合键失败: %v", err)
	}
	paymentBytes, err := ctx.GetStub().GetState(paymentKey)
	if err != nil {
		return nil, fmt.Errorf("[QueryPayment] 查询支付信息失败: %v", err)
	}
	if paymentBytes == nil {
		return nil, fmt.Errorf("[QueryPayment] 支付信息不存在: %s", paymentUUID)
	}

	var payment models.Payment
	err = json.Unmarshal(paymentBytes, &payment)
	if err != nil {
		return nil, fmt.Errorf("[QueryPayment] 解析支付信息失败: %v", err)
	}

	return &payment, nil
}

// PayForTransaction 支付房产交易（仅银行和投资者可调用）
func (s *SmartContract) PayForTransaction(ctx contractapi.TransactionContextInterface,
	transactionUUID string,
	paymentUUID string,
	paymentType string,
	amount float64,
	fromCitizenIDHash string,
	fromOrganization string,
	toCitizenIDHash string,
	toOrganization string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.BankMSP && clientMSPID != constances.InvestorMSP {
		return fmt.Errorf("[PayForTransaction] 只有银行和投资者可以支付交易")
	}

	// 检查交易是否已存在
	transactionKey, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{transactionUUID}...)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
	}
	transactionPublicBytes, err := ctx.GetStub().GetState(transactionKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询交易信息失败: %v", err)
	}
	if transactionPublicBytes == nil {
		return fmt.Errorf("[PayForTransaction] 交易不存在: %s", transactionUUID)
	}

	// 检查支付信息是否已存在
	paymentKey, err := s.createCompositeKey(ctx, constances.DocTypePayment, []string{paymentUUID}...)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
	}
	paymentBytes, err := ctx.GetStub().GetState(paymentKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询支付信息失败: %v", err)
	}
	if paymentBytes != nil {
		return fmt.Errorf("[PayForTransaction] 支付信息已存在: %s", paymentUUID)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 获取交易时间戳失败: %v", err)
	}

	// 创建支付信息
	payment := models.Payment{
		DocType:               constances.DocTypePayment,
		PaymentUUID:           paymentUUID,
		TransactionUUID:       transactionUUID,
		Amount:                amount,
		PaymentType:           paymentType,
		PayerCitizenIDHash:    fromCitizenIDHash,
		PayerOrganization:     fromOrganization,
		ReceiverCitizenIDHash: toCitizenIDHash,
		ReceiverOrganization:  toOrganization,
		CreateTime:            time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
	}

	// 序列化支付信息
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化支付信息失败: %v", err)
	}

	// 保存支付信息
	err = ctx.GetStub().PutState(paymentKey, paymentJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存支付信息失败: %v", err)
	}

	fromUserKey, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{fromCitizenIDHash, fromOrganization}...)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
	}
	// 查询fromCitizenIDHash的余额
	fromCitizenPrivateBytes, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, fromUserKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户余额失败: %v", err)
	}
	if fromCitizenPrivateBytes == nil {
		return fmt.Errorf("[PayForTransaction] 来源用户不存在: %s", fromCitizenIDHash)
	}

	var fromCitizenPrivate models.UserPrivate
	err = json.Unmarshal(fromCitizenPrivateBytes, &fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户余额失败: %v", err)
	}

	if fromCitizenPrivate.Balance < amount {
		return fmt.Errorf("[PayForTransaction] 余额不足: %s", fromCitizenIDHash)
	}

	// 更新fromCitizenIDHash的余额
	fromCitizenPrivate.Balance -= amount

	// 序列化fromCitizenIDHash的余额
	fromCitizenPrivateJSON, err := json.Marshal(fromCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户余额失败: %v", err)
	}

	// 保存fromCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(constances.UserDataCollection, fromUserKey, fromCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户余额失败: %v", err)
	}

	// 更新fromCitizenIDHash的
	fromCitizenPublicBytes, err := ctx.GetStub().GetState(fromUserKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户信息失败: %v", err)
	}

	var fromCitizenPublic models.UserPublic
	err = json.Unmarshal(fromCitizenPublicBytes, &fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户信息失败: %v", err)
	}

	fromCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化fromCitizenIDHash的
	fromCitizenPublicJSON, err := json.Marshal(fromCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户信息失败: %v", err)
	}

	// 保存fromCitizenIDHash的
	err = ctx.GetStub().PutState(fromUserKey, fromCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户信息失败: %v", err)
	}

	toUserKey, err := s.createCompositeKey(ctx, constances.DocTypeUser, []string{toCitizenIDHash, toOrganization}...)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 创建复合键失败: %v", err)
	}

	// 查询toCitizenIDHash的余额
	toCitizenBytes, err := ctx.GetStub().GetPrivateData(constances.UserDataCollection, toUserKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户余额失败: %v", err)
	}
	if toCitizenBytes == nil {
		return fmt.Errorf("[PayForTransaction] 目标用户不存在: %s", toCitizenIDHash)
	}

	var toCitizenPrivate models.UserPrivate
	err = json.Unmarshal(toCitizenBytes, &toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash的余额
	toCitizenPrivate.Balance += amount

	// 序列化toCitizenIDHash的余额
	toCitizenPrivateJSON, err := json.Marshal(toCitizenPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户余额失败: %v", err)
	}

	// 保存toCitizenIDHash的余额
	err = ctx.GetStub().PutPrivateData(constances.UserDataCollection, toUserKey, toCitizenPrivateJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户余额失败: %v", err)
	}

	// 更新toCitizenIDHash
	toCitizenPublicBytes, err := ctx.GetStub().GetState(toUserKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询用户信息失败: %v", err)
	}

	var toCitizenPublic models.UserPublic
	err = json.Unmarshal(toCitizenPublicBytes, &toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析用户信息失败: %v", err)
	}

	toCitizenPublic.LastUpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化toCitizenIDHash
	toCitizenPublicJSON, err := json.Marshal(toCitizenPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化用户信息失败: %v", err)
	}

	// 保存toCitizenIDHash
	err = ctx.GetStub().PutState(toUserKey, toCitizenPublicJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存用户信息失败: %v", err)
	}

	// 将该笔支付纳入交易
	transactionPrivateBytes, err := ctx.GetStub().GetPrivateData(constances.TransactionPrivateCollection, transactionKey)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 查询交易信息失败: %v", err)
	}

	var transactionPrivate models.TransactionPrivate
	err = json.Unmarshal(transactionPrivateBytes, &transactionPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析交易信息失败: %v", err)
	}

	transactionPrivate.PaymentUUIDList = append(transactionPrivate.PaymentUUIDList, paymentUUID)

	// 序列化交易
	transactionPrivateJSON, err := json.Marshal(transactionPrivate)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易
	err = ctx.GetStub().PutPrivateData(constances.TransactionPrivateCollection, transactionUUID, transactionPrivateJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存交易信息失败: %v", err)
	}

	// 更新交易状态
	var transactionPublic models.TransactionPublic
	err = json.Unmarshal(transactionPublicBytes, &transactionPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 解析交易信息失败: %v", err)
	}

	transactionPublic.Status = constances.TxStatusCompleted
	transactionPublic.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化交易
	transactionPublicJSON, err := json.Marshal(transactionPublic)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 序列化交易信息失败: %v", err)
	}

	// 保存交易
	err = ctx.GetStub().PutState(transactionUUID, transactionPublicJSON)
	if err != nil {
		return fmt.Errorf("[PayForTransaction] 保存交易信息失败: %v", err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 审计相关（待开发）
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// AuditTransaction 审计交易（仅审计机构可调用）
func (s *SmartContract) AuditTransaction(ctx contractapi.TransactionContextInterface, txID string, auditResult string, comments string) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.AuditMSP {
		return fmt.Errorf("只有审计机构可以进行交易审计")
	}

	// 查询交易信息
	for _, status := range []string{constances.TxStatusPending, constances.TxStatusCompleted} {
		txKey, err := s.createCompositeKey(ctx, constances.DocTypeTransaction, []string{status, txID}...)
		if err != nil {
			return err
		}

		transactionBytes, err := ctx.GetStub().GetState(txKey)
		if err != nil {
			return fmt.Errorf("查询交易信息失败: %v", err)
		}
		if transactionBytes != nil {
			// 创建审计记录
			auditRecord := AuditRecord{
				AuditID:      fmt.Sprintf("AUDIT_%s_%s", txID, time.Now().Format("20060102150405")),
				TargetType:   constances.DocTypeTransaction,
				TargetID:     txID,
				AuditorID:    clientMSPID,
				AuditorOrgID: constances.AuditMSP,
				Result:       auditResult,
				Comments:     comments,
				AuditedAt:    time.Now(),
			}

			// 创建审计记录复合键
			auditKey, err := s.createCompositeKey(ctx, constances.DocTypeAudit, []string{txID, auditRecord.AuditID}...)
			if err != nil {
				return err
			}

			// 序列化并保存审计记录
			auditJSON, err := json.Marshal(auditRecord)
			if err != nil {
				return fmt.Errorf("序列化审计记录失败: %v", err)
			}

			return ctx.GetStub().PutState(auditKey, auditJSON)
		}
	}

	return fmt.Errorf("交易ID %s 不存在", txID)
}

// QueryAuditHistory 查询交易的审计历史
func (s *SmartContract) QueryAuditHistory(ctx contractapi.TransactionContextInterface, txID string) ([]AuditRecord, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(constances.DocTypeAudit, []string{txID})
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %v", err)
	}
	defer iterator.Close()

	var records []AuditRecord
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条审计记录失败: %v", err)
		}

		var record AuditRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, fmt.Errorf("解析审计记录失败: %v", err)
		}

		records = append(records, record)
	}

	return records, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 其他功能
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// InitLedger 初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("初始化账本")
	return nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 合同相关
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateContract 创建合同（仅投资者和政府机构可以调用）
func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface,
	contractUUID string,
	docHash string,
	contractType string,
	creatorCitizenIDHash string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP {
		return fmt.Errorf("[CreateContract] 只有投资者和政府机构可以创建合同")
	}

	// 创建合同信息复合键
	key, err := s.createCompositeKey(ctx, constances.DocTypeContract, []string{contractUUID}...)
	if err != nil {
		return err
	}

	// 检查合同是否已存在
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("[CreateContract] 查询合同信息失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("[CreateContract] 合同UUID %s 已存在", contractUUID)
	}

	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[CreateContract] 获取交易时间戳失败: %v", err)
	}

	// 创建合同信息
	contract := models.Contract{
		ContractUUID:         contractUUID,
		CreateTime:           time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		UpdateTime:           time.Unix(now.Seconds, int64(now.Nanos)).UTC(),
		DocHash:              docHash,
		CreatorCitizenIDHash: creatorCitizenIDHash,
		ContractType:         contractType,
		Status:               constances.ContractStatusNormal,
	}

	// 序列化并保存合同信息
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("[CreateContract] 序列化合同信息失败: %v", err)
	}

	return ctx.GetStub().PutState(key, contractJSON)
}

// QueryContract 查询合同信息（仅投资者、政府机构和审计机构可以调用）
func (s *SmartContract) QueryContract(ctx contractapi.TransactionContextInterface,
	contractUUID string,
) (*models.Contract, error) {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return nil, err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP && clientMSPID != constances.AuditMSP {
		return nil, fmt.Errorf("[QueryContract] 只有投资者、政府机构和审计机构可以查询合同信息")
	}
	key, err := s.createCompositeKey(ctx, constances.DocTypeContract, []string{contractUUID}...)
	if err != nil {
		return nil, fmt.Errorf("[QueryContract] 创建复合键失败: %v", err)
	}
	contractBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("[QueryContract] 查询合同信息失败: %v", err)
	}
	if contractBytes == nil {
		return nil, fmt.Errorf("[QueryContract] 合同UUID %s 不存在", contractUUID)
	}

	var contract models.Contract
	err = json.Unmarshal(contractBytes, &contract)
	if err != nil {
		return nil, fmt.Errorf("[QueryContract] 解析合同信息失败: %v", err)
	}

	return &contract, nil
}

// UpdateContract 更新合同状态（仅投资者、政府机构、审计机构可以调用）
func (s *SmartContract) UpdateContract(ctx contractapi.TransactionContextInterface,
	contractUUID string,
	docHash string,
	contractType string,
	status string,
) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return err
	}

	if clientMSPID != constances.InvestorMSP && clientMSPID != constances.GovernmentMSP && clientMSPID != constances.AuditMSP {
		return fmt.Errorf("[UpdateContractStatus] 只有投资者、政府机构和审计机构可以更新合同状态")
	}

	// 检查合同是否存在
	contract, err := s.QueryContract(ctx, contractUUID)
	if err != nil {
		return fmt.Errorf("[UpdateContractStatus] 查询合同信息失败: %v", err)
	}

	// 检查合同状态
	if contract.Status == constances.ContractStatusFrozen {
		return fmt.Errorf("[UpdateContractStatus] 合同已冻结，无法更新状态")
	}

	// 更新合同
	if docHash != "" {
		contract.DocHash = docHash
	}
	if status != "" {
		contract.Status = status
	}
	if contractType != "" {
		contract.ContractType = contractType
	}
	now, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("[UpdateContractStatus] 获取交易时间戳失败: %v", err)
	}
	contract.UpdateTime = time.Unix(now.Seconds, int64(now.Nanos)).UTC()

	// 序列化并保存合同信息
	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("[UpdateContractStatus] 序列化合同信息失败: %v", err)
	}

	return ctx.GetStub().PutState(contract.ContractUUID, contractJSON)
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("创建智能合约失败: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动智能合约失败: %v", err)
	}
}
