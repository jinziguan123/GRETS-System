package chat_dto

import "time"

// ChatRoomDTO 聊天室DTO
type ChatRoomDTO struct {
	ID                  int64     `json:"id"`
	RoomUUID            string    `json:"roomUUID"`
	RealtyCertHash      string    `json:"realtyCertHash"`
	RealtyCert          string    `json:"realtyCert"`
	BuyerCitizenIDHash  string    `json:"buyerCitizenIDHash"`
	BuyerOrganization   string    `json:"buyerOrganization"`
	SellerCitizenIDHash string    `json:"sellerCitizenIDHash"`
	SellerOrganization  string    `json:"sellerOrganization"`
	Status              string    `json:"status"`
	VerificationAmount  float64   `json:"verificationAmount"`
	CreateTime          time.Time `json:"createTime"`
	CloseTime           time.Time `json:"closeTime"`
	LastMessageTime     time.Time `json:"lastMessageTime"`
	LastMessageContent  string    `json:"lastMessageContent"`
	UnreadCountBuyer    int       `json:"unreadCountBuyer"`
	UnreadCountSeller   int       `json:"unreadCountSeller"`
	// 扩展字段
	RealtyInfo      *RealtyInfoDTO      `json:"realtyInfo,omitempty"`
	ParticipantInfo *ParticipantInfoDTO `json:"participantInfo,omitempty"`
}

// ChatMessageDTO 聊天消息DTO
type ChatMessageDTO struct {
	ID                  int64     `json:"id"`
	MessageUUID         string    `json:"messageUUID"`
	RoomUUID            string    `json:"roomUUID"`
	SenderCitizenIDHash string    `json:"senderCitizenIDHash"`
	SenderOrganization  string    `json:"senderOrganization"`
	SenderName          string    `json:"senderName"`
	MessageType         string    `json:"messageType"`
	Content             string    `json:"content"`
	FileURL             string    `json:"fileURL"`
	FileName            string    `json:"fileName"`
	FileSize            int64     `json:"fileSize"`
	CreateTime          time.Time `json:"createTime"`
	IsRead              bool      `json:"isRead"`
	IsSelf              bool      `json:"isSelf"` // 是否为当前用户发送的消息
}

// VerificationRecordDTO 验资记录DTO
type VerificationRecordDTO struct {
	ID                 int64     `json:"id"`
	VerificationUUID   string    `json:"verificationUUID"`
	UserCitizenIDHash  string    `json:"userCitizenIDHash"`
	UserOrganization   string    `json:"userOrganization"`
	RealtyCertHash     string    `json:"realtyCertHash"`
	VerificationAmount float64   `json:"verificationAmount"`
	UserBalance        float64   `json:"userBalance"`
	Status             string    `json:"status"`
	CreateTime         time.Time `json:"createTime"`
}

// RealtyInfoDTO 房产信息DTO（用于聊天室中显示）
type RealtyInfoDTO struct {
	RealtyCert string   `json:"realtyCert"`
	RealtyType string   `json:"realtyType"`
	Price      float64  `json:"price"`
	Area       float64  `json:"area"`
	Address    string   `json:"address"`
	Images     []string `json:"images"`
}

// ParticipantInfoDTO 参与者信息DTO
type ParticipantInfoDTO struct {
	BuyerName  string `json:"buyerName"`
	SellerName string `json:"sellerName"`
}

// 请求DTO

// VerifyCapitalDTO 验资请求DTO
type VerifyCapitalDTO struct {
	UserCitizenID      string  `json:"userCitizenID" binding:"required"`
	UserOrganization   string  `json:"userOrganization" binding:"required"`
	RealtyCert         string  `json:"realtyCert" binding:"required"`
	VerificationAmount float64 `json:"verificationAmount" binding:"required,min=0"`
}

// CreateChatRoomDTO 创建聊天室请求DTO
type CreateChatRoomDTO struct {
	UserCitizenID      string  `json:"userCitizenID" binding:"required"`
	UserOrganization   string  `json:"userOrganization" binding:"required"`
	RealtyCert         string  `json:"realtyCert" binding:"required"`
	VerificationAmount float64 `json:"verificationAmount" binding:"required,min=0"`
}

// SendMessageDTO 发送消息请求DTO
type SendMessageDTO struct {
	UserCitizenID    string `json:"userCitizenID" binding:"required"`
	UserOrganization string `json:"userOrganization" binding:"required"`
	RoomUUID         string `json:"roomUUID" binding:"required"`
	MessageType      string `json:"messageType" binding:"required"`
	Content          string `json:"content" binding:"required"`
	FileURL          string `json:"fileURL"`
	FileName         string `json:"fileName"`
	FileSize         int64  `json:"fileSize"`
}

// QueryChatRoomListDTO 查询聊天室列表请求DTO
type QueryChatRoomListDTO struct {
	UserCitizenID    string `json:"userCitizenID" binding:"required"`
	UserOrganization string `json:"userOrganization" binding:"required"`
	Status           string `json:"status"`
	RealtyCert       string `json:"realtyCert"`
	PageSize         int    `json:"pageSize"`
	PageNumber       int    `json:"pageNumber"`
}

// QueryChatMessageListDTO 查询聊天消息列表请求DTO
type QueryChatMessageListDTO struct {
	UserCitizenID    string `json:"userCitizenID" binding:"required"`
	UserOrganization string `json:"userOrganization" binding:"required"`
	RoomUUID         string `json:"roomUUID" binding:"required"`
	MessageType      string `json:"messageType"`
	PageSize         int    `json:"pageSize"`
	PageNumber       int    `json:"pageNumber"`
}

// MarkMessagesReadDTO 标记消息已读请求DTO
type MarkMessagesReadDTO struct {
	UserCitizenID    string `json:"userCitizenID" binding:"required"`
	UserOrganization string `json:"userOrganization" binding:"required"`
	RoomUUID         string `json:"roomUUID" binding:"required"`
}

// CloseChatRoomDTO 关闭聊天室请求DTO
type CloseChatRoomDTO struct {
	UserCitizenID    string `json:"userCitizenID" binding:"required"`
	UserOrganization string `json:"userOrganization" binding:"required"`
	RoomUUID         string `json:"roomUUID" binding:"required"`
}

// 响应DTO

// VerifyCapitalResponseDTO 验资响应DTO
type VerifyCapitalResponseDTO struct {
	Success            bool    `json:"success"`
	VerificationUUID   string  `json:"verificationUUID"`
	UserBalance        float64 `json:"userBalance"`
	VerificationAmount float64 `json:"verificationAmount"`
	Message            string  `json:"message"`
}

// ChatRoomListResponseDTO 聊天室列表响应DTO
type ChatRoomListResponseDTO struct {
	ChatRooms  []ChatRoomDTO `json:"chatRooms"`
	Total      int           `json:"total"`
	PageSize   int           `json:"pageSize"`
	PageNumber int           `json:"pageNumber"`
}

// ChatMessageListResponseDTO 聊天消息列表响应DTO
type ChatMessageListResponseDTO struct {
	Messages   []ChatMessageDTO `json:"messages"`
	Total      int              `json:"total"`
	PageSize   int              `json:"pageSize"`
	PageNumber int              `json:"pageNumber"`
}
