import request from '@/utils/request'

// 获取用户的聊天室列表
export function getChatRoomList() {
  return request({
    url: '/api/v1/chatRooms',
    method: 'get'
  })
}

// 创建聊天室 (验资后针对特定房产创建聊天)
export function createChatRoom(data: {
  realtyCert: string
  toUserID: string
}) {
  return request({
    url: '/api/v1/chatRooms',
    method: 'post',
    data
  })
}

// 获取聊天消息
export function getChatMessages(roomId: string, params: {
  page?: number
  pageSize?: number
}) {
  return request({
    url: `/api/v1/chatRooms/${roomId}/messages`,
    method: 'get',
    params
  })
}

// 发送消息
export function sendMessage(roomId: string, data: {
  content: string
  messageType: 'text' | 'file'
  fileUrl?: string
}) {
  return request({
    url: `/api/v1/chatRooms/${roomId}/messages`,
    method: 'post',
    data
  })
}

// 上传聊天文件
export function uploadChatFile(roomId: string, file: File) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: `/api/v1/chatRooms/${roomId}/upload`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 验资请求 (购买二手房需要验资才能获取聊天权限)
export function verifyFunds(realtyCert: string) {
  return request({
    url: '/api/v1/funds/verify',
    method: 'post',
    data: { realtyCert }
  })
} 