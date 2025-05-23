import request from '@/utils/request'

// 类型定义
interface VerifyCapitalRequest {
  userCitizenID: string
  userOrganization: string
  realtyCert: string
  verificationAmount: number
}

interface CreateChatRoomRequest {
  userCitizenID: string
  userOrganization: string
  realtyCert: string
  verificationAmount: number
}

interface SendMessageRequest {
  userCitizenID: string
  userOrganization: string
  roomUUID: string
  messageType: string
  content: string
  fileURL?: string
  fileName?: string
  fileSize?: number
}

interface ChatRoomListRequest {
  userCitizenID: string
  userOrganization: string
  status?: string
  realtyCert?: string
  pageSize?: number
  pageNumber?: number
}

interface ChatMessageListRequest {
  userCitizenID: string
  userOrganization: string
  roomUUID: string
  messageType?: string
  pageSize?: number
  pageNumber?: number
}

interface MarkMessagesReadRequest {
  userCitizenID: string
  userOrganization: string
  roomUUID: string
}

interface CloseChatRoomRequest {
  userCitizenID: string
  userOrganization: string
  roomUUID: string
}

// 验资
export function verifyCapital(data: VerifyCapitalRequest) {
  return request({
    url: '/chat/verifyCapital',
    method: 'post',
    data
  })
}

// 创建聊天室
export function createChatRoom(data: CreateChatRoomRequest) {
  return request({
    url: '/chat/createChatRoom',
    method: 'post',
    data
  })
}

// 获取聊天室列表
export function getChatRoomList(data: ChatRoomListRequest) {
  return request({
    url: '/chat/getChatRoomList',
    method: 'post',
    data
  })
}

// 发送消息
export function sendMessage(data: SendMessageRequest) {
  return request({
    url: '/chat/sendMessage',
    method: 'post',
    data
  })
}

// 获取聊天消息列表
export function getChatMessageList(data: ChatMessageListRequest) {
  return request({
    url: '/chat/getChatMessageList',
    method: 'post',
    data
  })
}

// 标记消息为已读
export function markMessagesRead(data: MarkMessagesReadRequest) {
  return request({
    url: '/chat/markMessagesRead',
    method: 'post',
    data
  })
}

// 关闭聊天室
export function closeChatRoom(data: CloseChatRoomRequest) {
  return request({
    url: '/chat/closeChatRoom',
    method: 'post',
    data
  })
} 