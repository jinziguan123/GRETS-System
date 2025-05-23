<template>
  <div class="chat-room-container">
    <!-- 页面头部 -->
    <div class="chat-header">
      <div class="header-left">
        <el-button @click="goBack" circle>
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="room-info">
          <h3>{{ roomInfo.title }}</h3>
          <span class="room-subtitle">{{ roomInfo.subtitle }}</span>
        </div>
      </div>
      <div class="header-right">
        <el-tag :type="getRoomStatusType(roomInfo.status)">
          {{ getRoomStatusText(roomInfo.status) }}
        </el-tag>
      </div>
    </div>

    <!-- 聊天内容区域 -->
    <div class="chat-content" ref="chatContentRef">
      <div class="message-list">
        <div 
          v-for="message in messages" 
          :key="message.messageUUID" 
          :class="['message-item', message.isSelf ? 'self' : 'other']"
        >
          <div class="message-avatar">
            <el-avatar :size="40">
              {{ message.senderName ? message.senderName.charAt(0) : '?' }}
            </el-avatar>
          </div>
          <div class="message-content">
            <div class="message-header">
              <span class="sender-name">{{ message.senderName }}</span>
              <span class="message-time">{{ formatTime(message.createTime) }}</span>
            </div>
            <div class="message-body">
              <!-- 文本消息 -->
              <div v-if="message.messageType === 'TEXT'" class="text-message">
                {{ message.content }}
              </div>
              <!-- 系统消息 -->
              <div v-else-if="message.messageType === 'SYSTEM'" class="system-message">
                {{ message.content }}
              </div>
              <!-- 图片消息 -->
              <div v-else-if="message.messageType === 'IMAGE'" class="image-message">
                <el-image 
                  :src="message.fileURL" 
                  :preview-src-list="[message.fileURL]"
                  fit="cover"
                  style="width: 200px; height: 150px; border-radius: 8px;"
                />
              </div>
              <!-- 文件消息 -->
              <div v-else-if="message.messageType === 'FILE'" class="file-message">
                <el-button type="primary" link @click="downloadFile(message)">
                  <el-icon><Document /></el-icon>
                  {{ message.fileName }} ({{ formatFileSize(message.fileSize) }})
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 消息输入区域 -->
    <div class="chat-input" v-if="roomInfo.status === 'ACTIVE'">
      <div class="input-toolbar">
        <el-upload
          :action="uploadUrl"
          :show-file-list="false"
          :on-success="handleFileUpload"
          :before-upload="beforeFileUpload"
          accept="image/*,.pdf,.doc,.docx,.txt"
        >
          <el-button circle>
            <el-icon><Plus /></el-icon>
          </el-button>
        </el-upload>
      </div>
      <div class="input-content">
        <el-input
          v-model="inputMessage"
          type="textarea"
          :rows="3"
          placeholder="请输入消息..."
          @keydown.ctrl.enter="sendMessage"
          resize="none"
        />
      </div>
      <div class="input-actions">
        <el-button @click="clearInput">清空</el-button>
        <el-button type="primary" @click="sendMessage" :loading="sending">
          发送 (Ctrl+Enter)
        </el-button>
      </div>
    </div>

    <!-- 聊天室已关闭提示 -->
    <div v-else class="chat-closed">
      <el-empty description="聊天室已关闭或被冻结" />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Plus, Document } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getChatMessageList, sendMessage as sendChatMessage, markMessagesRead } from '@/api/chat'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 响应式数据
const chatContentRef = ref(null)
const messages = ref([])
const inputMessage = ref('')
const sending = ref(false)
const loading = ref(true)
const roomUUID = ref(route.params.roomUUID)

const roomInfo = reactive({
  roomUUID: '',
  title: '聊天室',
  subtitle: '',
  status: 'ACTIVE'
})

// 分页信息
const pagination = reactive({
  pageSize: 50,
  pageNumber: 1,
  total: 0,
  hasMore: true
})

// 上传配置
const uploadUrl = 'http://localhost:8080/api/v1/picture/upload'

// 轮询定时器
let pollingTimer = null

// 获取消息列表
const fetchMessages = async (isLoadMore = false) => {
  try {
    if (!isLoadMore) {
      loading.value = true
    }
    
    const requestData = {
      userCitizenID: userStore.user.citizenID,
      userOrganization: userStore.user.organization,
      roomUUID: roomUUID.value,
      pageSize: pagination.pageSize,
      pageNumber: pagination.pageNumber
    }
    
    const response = await getChatMessageList(requestData)
    
    if (isLoadMore) {
      // 加载更多消息时，添加到开头
      messages.value = [...response.messages.reverse(), ...messages.value]
    } else {
      // 首次加载或刷新时，直接设置
      messages.value = response.messages || []
    }
    
    pagination.total = response.total
    pagination.hasMore = pagination.pageNumber * pagination.pageSize < pagination.total
    
    // 滚动到底部（仅首次加载时）
    if (!isLoadMore) {
      await nextTick()
      scrollToBottom()
    }
    
    // 标记消息为已读
    await markMessagesRead({
      userCitizenID: userStore.user.citizenID,
      userOrganization: userStore.user.organization,
      roomUUID: roomUUID.value
    })
    
  } catch (error) {
    console.error('获取消息列表失败:', error)
    ElMessage.error('获取消息列表失败')
  } finally {
    loading.value = false
  }
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }
  
  sending.value = true
  try {
    const messageData = {
      userCitizenID: userStore.user.citizenID,
      userOrganization: userStore.user.organization,
      roomUUID: roomUUID.value,
      messageType: 'TEXT',
      content: inputMessage.value.trim()
    }
    
    const response = await sendChatMessage(messageData)
    
    // 添加新消息到列表
    messages.value.push(response)
    
    // 清空输入框
    inputMessage.value = ''
    
    // 滚动到底部
    await nextTick()
    scrollToBottom()
    
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error(error.response?.data?.message || '发送消息失败')
  } finally {
    sending.value = false
  }
}

// 处理文件上传
const beforeFileUpload = (file) => {
  const maxSize = 10 * 1024 * 1024 // 10MB
  if (file.size > maxSize) {
    ElMessage.error('文件大小不能超过10MB')
    return false
  }
  return true
}

const handleFileUpload = async (response, file) => {
  try {
    if (response.code === 200) {
      const messageData = {
        userCitizenID: userStore.user.citizenID,
        userOrganization: userStore.user.organization,
        roomUUID: roomUUID.value,
        messageType: file.type.startsWith('image/') ? 'IMAGE' : 'FILE',
        content: `发送了${file.type.startsWith('image/') ? '图片' : '文件'}: ${file.name}`,
        fileURL: response.data.url,
        fileName: file.name,
        fileSize: file.size
      }
      
      const messageResponse = await sendChatMessage(messageData)
      messages.value.push(messageResponse)
      
      await nextTick()
      scrollToBottom()
      
      ElMessage.success('文件发送成功')
    } else {
      ElMessage.error(response.message || '文件上传失败')
    }
  } catch (error) {
    console.error('发送文件消息失败:', error)
    ElMessage.error('发送文件消息失败')
  }
}

// 下载文件
const downloadFile = (message) => {
  const link = document.createElement('a')
  link.href = message.fileURL
  link.download = message.fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 清空输入
const clearInput = () => {
  inputMessage.value = ''
}

// 滚动到底部
const scrollToBottom = () => {
  if (chatContentRef.value) {
    chatContentRef.value.scrollTop = chatContentRef.value.scrollHeight
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 格式化时间
const formatTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) { // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) { // 24小时内
    return `${Math.floor(diff / 3600000)}小时前`
  } else {
    return date.toLocaleString()
  }
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1048576).toFixed(1) + ' MB'
}

// 获取房间状态类型
const getRoomStatusType = (status) => {
  const statusMap = {
    'ACTIVE': 'success',
    'CLOSED': 'info',
    'FROZEN': 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取房间状态文本
const getRoomStatusText = (status) => {
  const statusMap = {
    'ACTIVE': '进行中',
    'CLOSED': '已关闭',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || status
}

// 开始轮询获取新消息
const startPolling = () => {
  pollingTimer = setInterval(() => {
    // 获取最新消息（这里可以优化为只获取比当前最新消息更新的消息）
    if (roomInfo.status === 'ACTIVE') {
      fetchMessages()
    }
  }, 3000) // 每3秒轮询一次
}

// 停止轮询
const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

// 初始化房间信息
const initRoomInfo = () => {
  roomInfo.roomUUID = roomUUID.value
  roomInfo.title = `房产交易咨询`
  roomInfo.subtitle = `房产证号: ${route.query.realtyCert || ''}`
}

// 生命周期
onMounted(() => {
  initRoomInfo()
  fetchMessages()
  startPolling()
})

// 页面销毁时清理定时器
onBeforeUnmount(() => {
  stopPolling()
})

// 监听路由变化
watch(() => route.params.roomUUID, (newRoomUUID) => {
  if (newRoomUUID) {
    roomUUID.value = newRoomUUID
    initRoomInfo()
    fetchMessages()
  }
})
</script>

<style scoped>
.chat-room-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.room-info h3 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.room-subtitle {
  font-size: 14px;
  color: #909399;
}

.chat-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: #f5f7fa;
}

.message-list {
  max-width: 800px;
  margin: 0 auto;
}

.message-item {
  display: flex;
  margin-bottom: 20px;
  gap: 12px;
}

.message-item.self {
  flex-direction: row-reverse;
}

.message-item.self .message-content {
  align-items: flex-end;
}

.message-item.self .message-body {
  background-color: #409eff;
  color: white;
}

.message-content {
  display: flex;
  flex-direction: column;
  max-width: 60%;
  gap: 4px;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.message-item.self .message-header {
  flex-direction: row-reverse;
}

.sender-name {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.message-time {
  font-size: 12px;
  color: #909399;
}

.message-body {
  background-color: white;
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  word-wrap: break-word;
}

.system-message {
  background-color: #f0f9ff !important;
  color: #0ea5e9 !important;
  text-align: center;
  font-style: italic;
}

.file-message {
  padding: 8px 12px;
}

.chat-input {
  background: white;
  border-top: 1px solid #e4e7ed;
  padding: 16px 20px;
  display: flex;
  align-items: flex-end;
  gap: 12px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.input-toolbar {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-content {
  flex: 1;
}

.input-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.chat-closed {
  background: white;
  padding: 40px;
  text-align: center;
}

/* 滚动条样式 */
.chat-content::-webkit-scrollbar {
  width: 6px;
}

.chat-content::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.chat-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chat-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style> 