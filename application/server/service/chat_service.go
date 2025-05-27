package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"grets_server/dao"
	"grets_server/db"
	"grets_server/db/models"
	"grets_server/dto/chat_dto"
	chatDto "grets_server/dto/chat_dto"
	"grets_server/dto/realty_dto"
	"grets_server/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var GlobalChatService ChatService

func InitChatService() {
	GlobalChatService = NewChatService()
}

type chatService struct {
	db        *gorm.DB
	userDAO   *dao.UserDAO
	realtyDAO *dao.RealEstateDAO
}

type ChatService interface {
	// VerifyCapital 验资
	VerifyCapital(req *chatDto.VerifyCapitalDTO) (*chatDto.VerifyCapitalResponseDTO, error)
	// CreateChatRoom 创建聊天室
	CreateChatRoom(req *chatDto.CreateChatRoomDTO) (*chatDto.ChatRoomDTO, error)
	// SendMessage 发送消息
	SendMessage(req *chatDto.SendMessageDTO) (*chatDto.ChatMessageDTO, error)
	// GetChatRoomList 获取聊天室列表
	GetChatRoomList(req *chatDto.QueryChatRoomListDTO) (*chatDto.ChatRoomListResponseDTO, error)
	// GetChatMessageList 获取聊天消息列表
	GetChatMessageList(req *chatDto.QueryChatMessageListDTO) (*chatDto.ChatMessageListResponseDTO, error)
	// MarkMessagesRead 标记消息为已读
	MarkMessagesRead(req *chatDto.MarkMessagesReadDTO) error
	// CloseChatRoom 关闭聊天室
	CloseChatRoom(req *chatDto.CloseChatRoomDTO) error
	// hashString 计算字符串的SHA256哈希
	hashString(input string) string
	// convertChatRoomToDTO 转换聊天室实体为DTO
	convertChatRoomToDTO(room *models.ChatRoom) *chat_dto.ChatRoomDTO
	// convertMessageToDTO 转换消息实体为DTO
	convertMessageToDTO(message *models.ChatMessage, currentUserCitizenIDHash string) *chat_dto.ChatMessageDTO
	// fillChatRoomExtraInfo 填充聊天室的额外信息
	fillChatRoomExtraInfo(dto *chat_dto.ChatRoomDTO)
	// generateAddress 生成地址信息
	generateAddress(realty *realty_dto.RealtyDTO) string
}

func NewChatService() ChatService {
	return &chatService{
		db:        db.GlobalMysql,
		userDAO:   dao.NewUserDAO(),
		realtyDAO: dao.NewRealEstateDAO(),
	}
}

// VerifyCapital 验资
func (s *chatService) VerifyCapital(req *chatDto.VerifyCapitalDTO) (*chatDto.VerifyCapitalResponseDTO, error) {
	// 1. 检查房产是否存在
	realty, err := GlobalRealtyService.GetRealtyByRealtyCert(req.RealtyCert)
	if err != nil {
		return nil, fmt.Errorf("房产不存在: %v", err)
	}

	// 2. 检查房产状态是否允许交易
	if realty.Status != "PENDING_SALE" {
		return nil, fmt.Errorf("房产状态不允许交易")
	}

	// 3. 检查是否为房产所有者（自己不能和自己聊天）
	userCitizenIDHash := s.hashString(req.UserCitizenID)
	if realty.CurrentOwnerCitizenIDHash == userCitizenIDHash && realty.CurrentOwnerOrganization == req.UserOrganization {
		return nil, fmt.Errorf("不能与自己创建聊天室")
	}

	// 4. 检查余额是否充足
	balance, err := GlobalUserService.GetBalanceByCitizenIDAndOrganization(req.UserCitizenID, req.UserOrganization)
	if err != nil {
		return nil, fmt.Errorf("获取余额失败: %v", err)
	}

	if balance < req.VerificationAmount {
		return nil, fmt.Errorf("余额不足，当前余额: %.2f，需要验资: %.2f", balance, req.VerificationAmount)
	}

	// 5. 记录验资结果
	verificationUUID := uuid.New().String()

	return &chatDto.VerifyCapitalResponseDTO{
		Success:            true,
		VerificationUUID:   verificationUUID,
		UserBalance:        balance,
		VerificationAmount: req.VerificationAmount,
		Message:            "验资成功",
	}, nil
}

// CreateChatRoom 创建聊天室
func (s *chatService) CreateChatRoom(req *chatDto.CreateChatRoomDTO) (*chatDto.ChatRoomDTO, error) {
	// 2. 获取房产信息
	realty, err := GlobalRealtyService.GetRealtyByRealtyCertHash(utils.GenerateHash(req.RealtyCert))
	if err != nil {
		return nil, fmt.Errorf("房产不存在: %v", err)
	}

	userCitizenIDHash := s.hashString(req.UserCitizenID)
	realtyCertHash := s.hashString(req.RealtyCert)

	// 3. 检查是否已存在聊天室
	existingRoom := &models.ChatRoom{}
	err = s.db.Where("realty_cert_hash = ? AND buyer_citizen_id_hash = ? AND buyer_organization = ? AND status = 'ACTIVE'",
		realtyCertHash, userCitizenIDHash, req.UserOrganization).First(existingRoom).Error

	if err == nil {
		// 聊天室已存在，返回现有聊天室
		return s.convertChatRoomToDTO(existingRoom), nil
	}

	// 4. 创建新聊天室
	roomUUID := uuid.New().String()
	chatRoom := &models.ChatRoom{
		RoomUUID:            roomUUID,
		RealtyCertHash:      realtyCertHash,
		RealtyCert:          req.RealtyCert,
		BuyerCitizenIDHash:  userCitizenIDHash,
		BuyerOrganization:   req.UserOrganization,
		SellerCitizenIDHash: realty.CurrentOwnerCitizenIDHash,
		SellerOrganization:  realty.CurrentOwnerOrganization,
		Status:              "ACTIVE",
		VerificationAmount:  req.VerificationAmount,
		CreateTime:          time.Now(),
		LastMessageTime:     time.Now(),
		LastMessageContent:  "聊天室已创建",
		UnreadCountBuyer:    0,
		UnreadCountSeller:   1, // 卖方有一条系统消息未读
	}

	if err := s.db.Create(chatRoom).Error; err != nil {
		return nil, fmt.Errorf("创建聊天室失败: %v", err)
	}

	// 5. 发送系统欢迎消息
	welcomeMessage := &models.ChatMessage{
		MessageUUID:         uuid.New().String(),
		RoomUUID:            roomUUID,
		SenderCitizenIDHash: "system",
		SenderOrganization:  "system",
		SenderName:          "系统",
		MessageType:         "SYSTEM",
		Content:             fmt.Sprintf("欢迎进入房产 %s 的交易咨询聊天室！验资金额: %.2f 元", req.RealtyCert, req.VerificationAmount),
		CreateTime:          time.Now(),
		IsRead:              false,
	}

	if err := s.db.Create(welcomeMessage).Error; err != nil {
		return nil, fmt.Errorf("发送欢迎消息失败: %v", err)
	}

	return s.convertChatRoomToDTO(chatRoom), nil
}

// SendMessage 发送消息
func (s *chatService) SendMessage(req *chatDto.SendMessageDTO) (*chatDto.ChatMessageDTO, error) {
	// 1. 获取聊天室信息
	chatRoom := &models.ChatRoom{}
	err := s.db.Where("room_uuid = ?", req.RoomUUID).First(chatRoom).Error
	if err != nil {
		return nil, fmt.Errorf("聊天室不存在: %v", err)
	}

	// 2. 检查用户是否为聊天室参与者
	userCitizenIDHash := s.hashString(req.UserCitizenID)
	isBuyer := chatRoom.BuyerCitizenIDHash == userCitizenIDHash && chatRoom.BuyerOrganization == req.UserOrganization
	isSeller := chatRoom.SellerCitizenIDHash == userCitizenIDHash && chatRoom.SellerOrganization == req.UserOrganization

	if !isBuyer && !isSeller {
		return nil, fmt.Errorf("无权限在此聊天室发送消息")
	}

	// 3. 检查聊天室状态
	if chatRoom.Status != "ACTIVE" {
		return nil, fmt.Errorf("聊天室已关闭或被冻结")
	}

	// 4. 获取发送者信息
	user, err := s.userDAO.GetUserByCitizenID(req.UserCitizenID, req.UserOrganization)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 5. 创建消息
	message := &models.ChatMessage{
		MessageUUID:         uuid.New().String(),
		RoomUUID:            req.RoomUUID,
		SenderCitizenIDHash: userCitizenIDHash,
		SenderOrganization:  req.UserOrganization,
		SenderName:          user.Name,
		MessageType:         req.MessageType,
		Content:             req.Content,
		FileURL:             req.FileURL,
		FileName:            req.FileName,
		FileSize:            req.FileSize,
		CreateTime:          time.Now(),
		IsRead:              false,
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("发送消息失败: %v", err)
	}

	// 6. 更新聊天室的最后消息信息和未读计数
	updateData := map[string]interface{}{
		"last_message_time":    time.Now(),
		"last_message_content": req.Content,
	}

	if isBuyer {
		// 买方发送消息，卖方未读数+1
		updateData["unread_count_seller"] = gorm.Expr("unread_count_seller + ?", 1)
	} else {
		// 卖方发送消息，买方未读数+1
		updateData["unread_count_buyer"] = gorm.Expr("unread_count_buyer + ?", 1)
	}

	if err := s.db.Model(chatRoom).Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新聊天室信息失败: %v", err)
	}

	return s.convertMessageToDTO(message, userCitizenIDHash), nil
}

// GetChatRoomList 获取聊天室列表
func (s *chatService) GetChatRoomList(req *chatDto.QueryChatRoomListDTO) (*chatDto.ChatRoomListResponseDTO, error) {
	userCitizenIDHash := s.hashString(req.UserCitizenID)

	query := s.db.Model(&models.ChatRoom{}).Where(
		"(buyer_citizen_id_hash = ? AND buyer_organization = ?) OR (seller_citizen_id_hash = ? AND seller_organization = ?)",
		userCitizenIDHash, req.UserOrganization, userCitizenIDHash, req.UserOrganization,
	)

	// 添加筛选条件
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.RealtyCert != "" {
		query = query.Where("realty_cert = ?", req.RealtyCert)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询聊天室总数失败: %v", err)
	}

	// 分页查询
	var chatRooms []models.ChatRoom
	offset := (req.PageNumber - 1) * req.PageSize
	if err := query.Order("last_message_time DESC").Offset(offset).Limit(req.PageSize).Find(&chatRooms).Error; err != nil {
		return nil, fmt.Errorf("查询聊天室列表失败: %v", err)
	}

	// 转换为DTO
	chatRoomDTOs := make([]chatDto.ChatRoomDTO, 0, len(chatRooms))
	for _, room := range chatRooms {
		dto := s.convertChatRoomToDTO(&room)
		// 填充房产信息和参与者信息
		s.fillChatRoomExtraInfo(dto)
		chatRoomDTOs = append(chatRoomDTOs, *dto)
	}

	return &chatDto.ChatRoomListResponseDTO{
		ChatRooms:  chatRoomDTOs,
		Total:      int(total),
		PageSize:   req.PageSize,
		PageNumber: req.PageNumber,
	}, nil
}

// GetChatMessageList 获取聊天消息列表
func (s *chatService) GetChatMessageList(req *chatDto.QueryChatMessageListDTO) (*chatDto.ChatMessageListResponseDTO, error) {
	// 1. 检查用户是否有权限访问此聊天室
	userCitizenIDHash := s.hashString(req.UserCitizenID)
	chatRoom := &models.ChatRoom{}
	err := s.db.Where("room_uuid = ?", req.RoomUUID).First(chatRoom).Error
	if err != nil {
		return nil, fmt.Errorf("聊天室不存在: %v", err)
	}

	isBuyer := chatRoom.BuyerCitizenIDHash == userCitizenIDHash && chatRoom.BuyerOrganization == req.UserOrganization
	isSeller := chatRoom.SellerCitizenIDHash == userCitizenIDHash && chatRoom.SellerOrganization == req.UserOrganization

	if !isBuyer && !isSeller {
		return nil, fmt.Errorf("无权限访问此聊天室")
	}

	// 2. 查询消息
	query := s.db.Model(&models.ChatMessage{}).Where("room_uuid = ?", req.RoomUUID)

	if req.MessageType != "" {
		query = query.Where("message_type = ?", req.MessageType)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询消息总数失败: %v", err)
	}

	// 分页查询（按时间升序）
	var messages []models.ChatMessage
	offset := (req.PageNumber - 1) * req.PageSize
	if err := query.Order("create_time ASC").Offset(offset).Limit(req.PageSize).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %v", err)
	}

	// 转换为DTO
	messageDTOs := make([]chatDto.ChatMessageDTO, 0, len(messages))
	for _, message := range messages {
		dto := s.convertMessageToDTO(&message, userCitizenIDHash)
		messageDTOs = append(messageDTOs, *dto)
	}

	return &chatDto.ChatMessageListResponseDTO{
		Messages:   messageDTOs,
		Total:      int(total),
		PageSize:   req.PageSize,
		PageNumber: req.PageNumber,
	}, nil
}

// MarkMessagesRead 标记消息为已读
func (s *chatService) MarkMessagesRead(req *chatDto.MarkMessagesReadDTO) error {
	// 1. 检查用户权限
	userCitizenIDHash := s.hashString(req.UserCitizenID)
	chatRoom := &models.ChatRoom{}
	err := s.db.Where("room_uuid = ?", req.RoomUUID).First(chatRoom).Error
	if err != nil {
		return fmt.Errorf("聊天室不存在: %v", err)
	}

	isBuyer := chatRoom.BuyerCitizenIDHash == userCitizenIDHash && chatRoom.BuyerOrganization == req.UserOrganization
	isSeller := chatRoom.SellerCitizenIDHash == userCitizenIDHash && chatRoom.SellerOrganization == req.UserOrganization

	if !isBuyer && !isSeller {
		return fmt.Errorf("无权限访问此聊天室")
	}

	// 2. 标记消息为已读（标记对方发送的消息）
	updateQuery := s.db.Model(&models.ChatMessage{}).Where("room_uuid = ? AND is_read = false", req.RoomUUID)

	if isBuyer {
		// 买方读消息，标记卖方发送的消息为已读
		updateQuery = updateQuery.Where("sender_citizen_id_hash = ? AND sender_organization = ?",
			chatRoom.SellerCitizenIDHash, chatRoom.SellerOrganization)
	} else {
		// 卖方读消息，标记买方发送的消息为已读
		updateQuery = updateQuery.Where("sender_citizen_id_hash = ? AND sender_organization = ?",
			chatRoom.BuyerCitizenIDHash, chatRoom.BuyerOrganization)
	}

	if err := updateQuery.Update("is_read", true).Error; err != nil {
		return fmt.Errorf("标记消息已读失败: %v", err)
	}

	// 3. 重置聊天室未读计数
	updateData := map[string]interface{}{}
	if isBuyer {
		updateData["unread_count_buyer"] = 0
	} else {
		updateData["unread_count_seller"] = 0
	}

	if err := s.db.Model(chatRoom).Updates(updateData).Error; err != nil {
		return fmt.Errorf("重置未读计数失败: %v", err)
	}

	return nil
}

// CloseChatRoom 关闭聊天室
func (s *chatService) CloseChatRoom(req *chatDto.CloseChatRoomDTO) error {
	// 检查用户权限（只有卖方可以关闭聊天室）
	userCitizenIDHash := s.hashString(req.UserCitizenID)
	chatRoom := &models.ChatRoom{}
	err := s.db.Where("room_uuid = ?", req.RoomUUID).First(chatRoom).Error
	if err != nil {
		return fmt.Errorf("聊天室不存在: %v", err)
	}

	isSeller := chatRoom.SellerCitizenIDHash == userCitizenIDHash && chatRoom.SellerOrganization == req.UserOrganization
	if !isSeller {
		return fmt.Errorf("只有卖方可以关闭聊天室")
	}

	// 更新聊天室状态
	updateData := map[string]interface{}{
		"status":     "CLOSED",
		"close_time": time.Now(),
	}

	if err := s.db.Model(chatRoom).Updates(updateData).Error; err != nil {
		return fmt.Errorf("关闭聊天室失败: %v", err)
	}

	return nil
}

// 工具方法

// hashString 计算字符串的SHA256哈希
func (s *chatService) hashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// convertChatRoomToDTO 转换聊天室实体为DTO
func (s *chatService) convertChatRoomToDTO(room *models.ChatRoom) *chat_dto.ChatRoomDTO {
	dto := &chat_dto.ChatRoomDTO{
		ID:                  room.ID,
		RoomUUID:            room.RoomUUID,
		RealtyCertHash:      room.RealtyCertHash,
		RealtyCert:          room.RealtyCert,
		BuyerCitizenIDHash:  room.BuyerCitizenIDHash,
		BuyerOrganization:   room.BuyerOrganization,
		SellerCitizenIDHash: room.SellerCitizenIDHash,
		SellerOrganization:  room.SellerOrganization,
		Status:              room.Status,
		VerificationAmount:  room.VerificationAmount,
		CreateTime:          room.CreateTime,
		LastMessageTime:     room.LastMessageTime,
		LastMessageContent:  room.LastMessageContent,
		UnreadCountBuyer:    room.UnreadCountBuyer,
		UnreadCountSeller:   room.UnreadCountSeller,
	}

	if room.CloseTime.Valid {
		dto.CloseTime = room.CloseTime.Time
	}

	return dto
}

// convertMessageToDTO 转换消息实体为DTO
func (s *chatService) convertMessageToDTO(message *models.ChatMessage, currentUserCitizenIDHash string) *chat_dto.ChatMessageDTO {
	return &chat_dto.ChatMessageDTO{
		ID:                  message.ID,
		MessageUUID:         message.MessageUUID,
		RoomUUID:            message.RoomUUID,
		SenderCitizenIDHash: message.SenderCitizenIDHash,
		SenderOrganization:  message.SenderOrganization,
		SenderName:          message.SenderName,
		MessageType:         message.MessageType,
		Content:             message.Content,
		FileURL:             message.FileURL,
		FileName:            message.FileName,
		FileSize:            message.FileSize,
		CreateTime:          message.CreateTime,
		IsRead:              message.IsRead,
		IsSelf:              message.SenderCitizenIDHash == currentUserCitizenIDHash,
	}
}

// fillChatRoomExtraInfo 填充聊天室的额外信息
func (s *chatService) fillChatRoomExtraInfo(dto *chat_dto.ChatRoomDTO) {
	// 填充房产信息
	realty, err := GlobalRealtyService.GetRealtyByRealtyCert(dto.RealtyCert)
	if err == nil {
		dto.RealtyInfo = &chat_dto.RealtyInfoDTO{
			RealtyCert: realty.RealtyCert,
			RealtyType: realty.RealtyType,
			Price:      realty.Price,
			Area:       realty.Area,
			Address:    s.generateAddress(realty),
			Images:     realty.Images,
		}
	}

	// 填充参与者信息
	buyerUser, _ := GlobalUserService.GetUserByCitizenIDAndOrganization(dto.BuyerCitizenIDHash, dto.BuyerOrganization)
	sellerUser, _ := GlobalUserService.GetUserByCitizenIDAndOrganization(dto.SellerCitizenIDHash, dto.SellerOrganization)

	dto.ParticipantInfo = &chat_dto.ParticipantInfoDTO{}
	if buyerUser != nil {
		dto.ParticipantInfo.BuyerName = buyerUser.Name
	}
	if sellerUser != nil {
		dto.ParticipantInfo.SellerName = sellerUser.Name
	}
}

// generateAddress 生成地址信息
func (s *chatService) generateAddress(realty *realty_dto.RealtyDTO) string {
	if realty.Province == "" {
		return "地址不详"
	}

	address := realty.Province
	if realty.City != "" && realty.City != realty.Province {
		address += realty.City
	}
	if realty.District != "" {
		address += realty.District
	}
	if realty.Street != "" {
		address += realty.Street
	}
	if realty.Community != "" {
		address += realty.Community
	}

	return address
}
