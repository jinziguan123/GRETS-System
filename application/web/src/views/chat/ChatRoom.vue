<template>
  <div class="chat-room-wrapper">
    <!-- 聊天室头部 -->
    <div class="chat-header">
      <div class="header-left">
        <el-button 
          type="text" 
          @click="goBack" 
          class="back-btn"
          :icon="ArrowLeft"
        >
        </el-button>
        <div class="room-info">
          <h3 class="room-title">{{ roomInfo.title }}</h3>
          <span class="room-subtitle">{{ roomInfo.subtitle }}</span>
        </div>
      </div>
      <div class="header-right">
        <!-- 操作按钮 -->
        <div class="header-actions">
          <el-button 
            type="primary" 
            size="small" 
            @click="goToRealtyDetail"
            class="action-btn"
          >
            房产详情
          </el-button>
          <el-button 
            v-if="isBuyer || !isSeller" 
            type="success" 
            size="small" 
            @click="applyTransaction"
            class="action-btn"
            :disabled="roomInfo.status !== 'ACTIVE'"
          >
            申请交易
          </el-button>
        </div>
        
        <div class="connection-status" :class="{ 'connected': isConnected }">
          <span class="status-dot"></span>
          <span class="status-text">{{ isConnected ? '在线' : '连接中...' }}</span>
        </div>
        <el-tag 
          :type="getRoomStatusType(roomInfo.status)" 
          class="room-status-tag"
        >
          {{ getRoomStatusText(roomInfo.status) }}
        </el-tag>
      </div>
    </div>

    <!-- 聊天消息区域 -->
    <div class="chat-content" ref="chatContentRef">
      <div class="message-container">
        <div 
          v-for="message in messages" 
          :key="message.messageUUID" 
          :class="['message-wrapper', { 'message-self': message.isSelf }]"
        >
          <!-- 消息时间（间隔显示） -->
          <div class="message-time" v-if="shouldShowTime(message)">
            {{ formatTime(message.createTime) }}
          </div>
          
          <div class="message-item">
            <!-- 头像 -->
            <div class="message-avatar">
              <el-avatar 
                :size="36" 
                :style="{ backgroundColor: message.isSelf ? '#07c160' : '#4a90e2' }"
              >
                {{ getAvatarText(message) }}
              </el-avatar>
            </div>
            
            <!-- 消息内容 -->
            <div class="message-content">
              <!-- 发送者名称 -->
              <div v-if="!message.isSelf" class="sender-name">
                {{ message.senderName }}
              </div>
              
              <!-- 消息气泡 -->
              <div class="message-bubble" :class="getBubbleClass(message)">
                <!-- 文本消息 -->
                <div v-if="message.messageType === 'TEXT'" class="text-message">
                  {{ message.content }}
                </div>
                
                <!-- 系统消息 -->
                <div v-else-if="message.messageType === 'SYSTEM'" class="system-message">
                  <el-icon class="system-icon"><InfoFilled /></el-icon>
                  {{ message.content }}
                </div>
                
                <!-- 图片消息 -->
                <div v-else-if="message.messageType === 'IMAGE'" class="image-message">
                  <el-image 
                    :src="message.fileURL" 
                    :preview-src-list="[message.fileURL]"
                    fit="cover"
                    class="chat-image"
                    :preview-teleported="true"
                  />
                </div>
                
                <!-- 文件消息 -->
                <div v-else-if="message.messageType === 'FILE'" class="file-message">
                  <div class="file-icon">
                    <el-icon size="24"><Document /></el-icon>
                  </div>
                  <div class="file-info">
                    <div class="file-name">{{ message.fileName }}</div>
                    <div class="file-size">{{ formatFileSize(message.fileSize) }}</div>
                  </div>
                  <el-button 
                    type="primary" 
                    text 
                    @click="downloadFile(message)"
                    class="download-btn"
                  >
                    下载
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 正在输入提示 -->
        <div v-if="sending" class="typing-indicator">
          <div class="typing-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <span class="typing-text">发送中...</span>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="chat-input-area" v-if="roomInfo.status === 'ACTIVE'">
      <!-- 工具栏 -->
      <!-- <div class="input-toolbar">
        <el-upload
          :action="uploadUrl"
          :show-file-list="false"
          :on-success="handleFileUpload"
          :before-upload="beforeFileUpload"
          accept="image/*,.pdf,.doc,.docx,.txt"
          class="upload-btn"
        >
          <el-button 
            type="text" 
            class="toolbar-btn"
            :icon="Plus"
          >
          </el-button>
        </el-upload>
        
        <el-button 
          type="text" 
          class="toolbar-btn"
          :icon="Picture"
          @click="triggerImageUpload"
        >
        </el-button>
      </div> -->
      
      <!-- 输入框 -->
      <div class="input-container">
        <el-input
          v-model="inputMessage"
          type="textarea"
          :rows="3"
          resize="none"
          placeholder="请输入消息... (按 Ctrl+Enter 发送)"
          @keydown.ctrl.enter="sendMessage"
          @keydown.enter.prevent="handleEnterKey"
          class="message-input"
          :disabled="!isConnected"
        />
        
        <!-- 发送按钮 -->
        <div class="send-container">
          <el-button
            type="primary"
            @click="sendMessage"
            :loading="sending"
            :disabled="!inputMessage.trim() || !isConnected"
            @keydown.enter="handleEnter"
            class="send-btn"
          >
            发送
          </el-button>
        </div>
      </div>
    </div>

    <!-- 聊天室已关闭提示 -->
    <div v-else class="chat-disabled">
      <el-empty 
        description="聊天室已关闭或被冻结" 
        :image-size="100"
      >
        <template #image>
          <el-icon size="100" color="#c0c4cc"><ChatDotSquare /></el-icon>
        </template>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick, watch, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus, Document, InfoFilled, Picture, ChatDotSquare } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getChatMessageList, sendMessage as sendChatMessage, markMessagesRead } from '@/api/chat'
import CryptoJS, { SHA256 } from 'crypto-js'

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
const webSocket = ref(null)
const isConnected = ref(false)

const roomInfo = reactive({
  roomUUID: '',
  title: '聊天室',
  subtitle: '',
  status: 'ACTIVE',
  buyerCitizenIDHash: '',
  sellerCitizenIDHash: '',
  realtyCert: ''
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

// 判断消息是否为自己发送的
const isMessageFromSelf = (message) => {
  const currentUserCitizenIDHash = userStore.user.citizenIDHash || 
    CryptoJS.SHA256(userStore.user.citizenID).toString()
  
  return message.senderCitizenIDHash === currentUserCitizenIDHash && 
         message.senderOrganization === userStore.user.organization
}

// 判断当前用户是否为买方
const isBuyer = computed(() => {
  const currentUserCitizenIDHash = userStore.user.citizenIDHash || 
    CryptoJS.SHA256(userStore.user.citizenID).toString()
  
  return route.query.buyerCitizenIDHash === currentUserCitizenIDHash &&
         route.query.buyerOrganization === userStore.user.organization
})

const isSeller = computed(() => {
  const currentUserCitizenIDHash = userStore.user.citizenIDHash || 
    CryptoJS.SHA256(userStore.user.citizenID).toString()
  
  return route.query.sellerCitizenIDHash === currentUserCitizenIDHash &&
         route.query.sellerOrganization === userStore.user.organization
})

const handleEnter = (e) => {
  if (e.keyCode === 13 || e.keyCode === 108) {
    sendMessage()
  }
}

// WebSocket连接
const connectWebSocket = () => {
  const token = localStorage.getItem('token')
  // 通过查询参数传递token
  const wsUrl = `ws://localhost:8080/api/v1/chat/ws/${roomUUID.value}?token=${encodeURIComponent(token)}`
  
  webSocket.value = new WebSocket(wsUrl)
  
  webSocket.value.onopen = () => {
    console.log('WebSocket连接已建立')
    isConnected.value = true
  }
  
  webSocket.value.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)
      handleWebSocketMessage(message)
    } catch (error) {
      console.error('解析WebSocket消息失败:', error)
    }
  }
  
  webSocket.value.onclose = () => {
    console.log('WebSocket连接已关闭')
    isConnected.value = false
    
    // 重连逻辑
    setTimeout(() => {
      if (roomInfo.status === 'ACTIVE') {
        connectWebSocket()
      }
    }, 3000)
  }
  
  webSocket.value.onerror = (error) => {
    console.error('WebSocket错误:', error)
    isConnected.value = false
  }
}

// 处理WebSocket消息
const handleWebSocketMessage = (message) => {
  switch (message.type) {
    case 'join':
      console.log('成功加入聊天室')
      break
    case 'newMessage':
      // 收到新消息，添加到消息列表
      if (message.data) {
        const newMessage = {
          ...message.data,
          isSelf: isMessageFromSelf(message.data)
        }
        messages.value.push(newMessage)
        nextTick(() => {
          scrollToBottom()
        })
      }
      break
    case 'pong':
      // 心跳回应
      break
  }
}

// 断开WebSocket连接
const disconnectWebSocket = () => {
  if (webSocket.value) {
    webSocket.value.close()
    webSocket.value = null
    isConnected.value = false
  }
}

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
    
    // 为每条消息添加isSelf字段
    const messagesWithSelfFlag = (response.messages || []).map(message => ({
      ...message,
      isSelf: isMessageFromSelf(message)
    }))
    
    if (isLoadMore) {
      // 加载更多消息时，添加到开头
      messages.value = [...messagesWithSelfFlag.reverse(), ...messages.value]
    } else {
      // 首次加载或刷新时，直接设置
      messages.value = messagesWithSelfFlag
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
  
  if (!isConnected.value) {
    ElMessage.error('WebSocket连接已断开，请稍后重试')
    return
  }
  
  sending.value = true
  try {
    // 通过HTTP API发送消息到服务器
    const messageData = {
      userCitizenID: userStore.user.citizenID,
      userOrganization: userStore.user.organization,
      roomUUID: roomUUID.value,
      messageType: 'TEXT',
      content: inputMessage.value.trim()
    }
    
    const response = await sendChatMessage(messageData)
    
    // 消息发送成功，WebSocket会收到广播消息并自动添加到列表中
    // 清空输入框
    inputMessage.value = ''
    
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

// 获取头像显示文本
const getAvatarText = (message) => {
  if (message.messageType === 'SYSTEM') {
    return '系'
  }
  if (message.senderName) {
    return message.senderName.charAt(0).toUpperCase()
  }
  return '?'
}

// 获取消息气泡样式类
const getBubbleClass = (message) => {
  const classes = []
  
  if (message.messageType === 'SYSTEM') {
    classes.push('system-bubble')
  } else if (message.isSelf) {
    classes.push('self-bubble')
  } else {
    classes.push('other-bubble')
  }
  
  return classes
}

// 判断是否显示时间
const shouldShowTime = (message) => {
  const currentIndex = messages.value.findIndex(msg => msg.messageUUID === message.messageUUID)
  
  // 第一条消息总是显示时间
  if (currentIndex === 0) {
    return true
  }
  
  // 系统消息总是显示时间
  if (message.messageType === 'SYSTEM') {
    return true
  }
  
  // 获取上一条消息
  const previousMessage = messages.value[currentIndex - 1]
  if (!previousMessage || !previousMessage.createTime) {
    return true
  }
  
  // 计算时间差（毫秒）
  const currentTime = new Date(message.createTime).getTime()
  const previousTime = new Date(previousMessage.createTime).getTime()
  const timeDiff = currentTime - previousTime
  
  // 如果间隔超过10分钟（600000毫秒），显示时间
  return timeDiff > 600000
}

// 处理Enter键
const handleEnterKey = (event) => {
  if (!event.ctrlKey) {
    // 普通Enter键不发送，只换行
    return
  }
}

// 触发图片上传
const triggerImageUpload = () => {
  // 创建隐藏的文件输入框
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = (e) => {
    const file = e.target.files[0]
    if (file) {
      // 这里可以添加图片上传逻辑
      console.log('选择的图片:', file)
    }
  }
  input.click()
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

// 跳转到房产详情页面
const goToRealtyDetail = async () => {
  const realtyCertHash = route.query.realtyCertHash || roomInfo.realtyCertHash
  if (realtyCertHash) {
    router.push({
      path: `/realty/${realtyCertHash}`
    })
  } else {
    ElMessage.warning('房产证号不存在')
  }
}

// 申请交易
const applyTransaction = async () => {
  router.push(`/transaction/create?realtyCert=${route.query.realtyCert}`)
}

// 初始化房间信息
const initRoomInfo = () => {
  roomInfo.roomUUID = roomUUID.value
  roomInfo.subtitle = `房产交易咨询`
  roomInfo.title = `房产证号: ${route.query.realtyCert || ''}`
  roomInfo.realtyCert = route.query.realtyCert || ''
}

// 生命周期
onMounted(() => {
  initRoomInfo()
  fetchMessages()
  connectWebSocket()
  window.addEventListener('keydown', handleEnter);
})

// 页面销毁时清理WebSocket连接
onBeforeUnmount(() => {
  disconnectWebSocket()
})

// 监听路由变化
watch(() => route.params.roomUUID, (newRoomUUID) => {
  if (newRoomUUID) {
    roomUUID.value = newRoomUUID
    initRoomInfo()
    fetchMessages()
    disconnectWebSocket()
    connectWebSocket()
  }
})
</script>

<style scoped>
.chat-room-wrapper {
  height: 80vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  overflow-y: hidden; /* 防止整体页面滚动 */
}

/* 聊天室头部样式 */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  min-height: 60px; /* 确保头部高度固定 */
  padding: 0 20px;
  background: #ffffff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  position: relative;
  z-index: 100;
  flex-shrink: 0; /* 防止被压缩 */
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  padding: 8px;
  font-size: 18px;
  color: #606266;
  transition: color 0.2s;
}

.back-btn:hover {
  color: #409eff;
}

.room-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.room-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  line-height: 1.2;
}

.room-subtitle {
  font-size: 12px;
  color: #909399;
  line-height: 1;
}

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  .header-actions {
    gap: 4px;
  }
  
  .action-btn {
    font-size: 11px;
    padding: 4px 8px;
  }

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  font-size: 12px;
  padding: 6px 12px;
  border-radius: 4px;
  transition: all 0.2s;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #f56c6c;
  transition: background-color 0.3s;
}

.connection-status.connected .status-dot {
  background-color: #67c23a;
}

.status-text {
  font-size: 12px;
  color: #909399;
}

.room-status-tag {
  font-size: 12px;
}

/* 聊天消息区域样式 */
.chat-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 8px;
  background: linear-gradient(to bottom, #f5f5f5 0%, #f0f0f0 100%);
  min-height: 0; /* 确保可以被压缩 */
}

.message-container {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 8px;
}

.message-wrapper {
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
}

.message-wrapper.message-self {
  align-items: flex-end;
}

.message-time {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin: 12px 0 8px 0;
  background: rgba(0, 0, 0, 0.05);
  padding: 6px 16px;
  border-radius: 12px;
  display: inline-block;
  align-self: center;
}

.message-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  max-width: 85%;
}

.message-wrapper.message-self .message-item {
  flex-direction: row-reverse;
}

.message-avatar {
  flex: 0 0 auto;
  margin-top: 4px;
}

.message-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
  align-items: flex-start; /* 让气泡根据内容宽度自适应 */
}

.message-wrapper.message-self .message-content {
  align-items: flex-end;
}

.sender-name {
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
  padding: 0 12px;
}

/* 消息气泡样式 */
.message-bubble {
  position: relative;
  border-radius: 12px;
  padding: 14px 18px;
  word-wrap: break-word;
  word-break: break-word;
  min-width: 40px; /* 最小宽度 */
  max-width: 100%; /* 最大宽度 */
  width: fit-content; /* 根据内容自适应宽度 */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
  font-size: 15px;
  line-height: 1.6;
  display: inline-block; /* 确保宽度自适应 */
}

.message-bubble:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 对方消息气泡 */
.other-bubble {
  background: #ffffff;
  color: #333;
}

.other-bubble::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 16px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 8px 8px 8px 0;
  border-color: transparent #ffffff transparent transparent;
}

/* 自己发送的消息气泡 */
.self-bubble {
  background: #07c160;
  color: #ffffff;
}

.self-bubble::after {
  content: '';
  position: absolute;
  right: -8px;
  top: 16px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 8px 0 8px 8px;
  border-color: transparent transparent transparent #07c160;
}

/* 系统消息气泡 */
.system-bubble {
  background: #f0f9ff;
  color: #0ea5e9;
  border: 1px solid #e0f2fe;
  text-align: center;
  font-size: 12px;
  align-self: center;
  margin: 8px 0;
  width: auto; /* 自适应宽度 */
  display: inline-block;
}

.system-bubble::before,
.system-bubble::after {
  display: none;
}

.system-message {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-style: italic;
  white-space: nowrap; /* 防止系统消息换行 */
}

.system-icon {
  font-size: 14px;
}

/* 图片消息样式 */
.image-message {
  padding: 0;
  border-radius: 8px;
  overflow: hidden;
  width: auto; /* 自适应宽度 */
  display: inline-block;
}

.chat-image {
  max-width: 250px;
  max-height: 200px;
  width: auto; /* 保持图片原始比例 */
  height: auto;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
  display: block;
}

.chat-image:hover {
  transform: scale(1.02);
}

/* 文件消息样式 */
.file-message {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.9);
  padding: 12px;
  border-radius: 8px;
  border: 1px dashed #ddd;
  transition: all 0.2s;
  width: auto; /* 自适应宽度 */
  min-width: 200px; /* 文件消息最小宽度 */
}

.file-message:hover {
  border-color: #409eff;
  background: rgba(64, 158, 255, 0.05);
}

.file-icon {
  color: #409eff;
  font-size: 24px;
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.file-size {
  font-size: 12px;
  color: #999;
}

.download-btn {
  font-size: 12px;
  padding: 4px 8px;
}

/* 正在输入提示 */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  margin: 16px 0;
  align-self: flex-start;
  font-size: 12px;
  color: #666;
}

.typing-dots {
  display: flex;
  gap: 3px;
}

.typing-dots span {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background-color: #666;
  animation: typing-bounce 1.4s infinite ease-in-out;
}

.typing-dots span:nth-child(1) { animation-delay: -0.32s; }
.typing-dots span:nth-child(2) { animation-delay: -0.16s; }

@keyframes typing-bounce {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

/* 输入区域样式 */
.chat-input-area {
  background: #ffffff;
  border-top: 1px solid #e4e7ed;
  padding: 16px 20px;
  position: relative;
  z-index: 100;
  flex-shrink: 0; /* 防止被压缩 */
  min-height: auto; /* 确保高度自适应内容 */
}

.input-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.toolbar-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  color: #666;
  font-size: 16px;
  transition: all 0.2s;
  background: #f5f5f5;
}

.toolbar-btn:hover {
  background: #e6f7ff;
  color: #409eff;
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.message-input {
  flex: 1;
}

.message-input :deep(.el-textarea__inner) {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  padding: 12px;
  font-size: 14px;
  line-height: 1.4;
  resize: none;
  transition: border-color 0.2s;
  background: #fafafa;
}

.message-input :deep(.el-textarea__inner):focus {
  border-color: #409eff;
  background: #ffffff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.message-input :deep(.el-textarea__inner):disabled {
  background: #f5f5f5;
  color: #999;
}

.send-container {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.send-btn {
  height: 36px;
  padding: 0 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.send-btn:not(:disabled) {
  background: #07c160;
  border-color: #07c160;
}

.send-btn:not(:disabled):hover {
  background: #06ad56;
  border-color: #06ad56;
  transform: translateY(-1px);
}

.send-btn:disabled {
  background: #f5f5f5;
  border-color: #ddd;
  color: #999;
}

/* 聊天室已关闭样式 */
.chat-disabled {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  padding: 40px;
  min-height: 0; /* 确保可以被压缩 */
}

.chat-disabled :deep(.el-empty__description) {
  color: #909399;
  font-size: 14px;
}

/* 滚动条样式 */
.chat-content::-webkit-scrollbar {
  width: 6px;
}

.chat-content::-webkit-scrollbar-track {
  background: transparent;
}

.chat-content::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
  transition: background 0.2s;
}

.chat-content::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.4);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-room-wrapper {
    height: 100vh;
    overflow: hidden; /* 确保移动端也不可滚动 */
  }
  
  .chat-header {
    height: 50px;
    min-height: 50px;
    padding: 0 16px;
    flex-shrink: 0;
  }
  
  .room-title {
    font-size: 14px;
  }
  
  .room-subtitle {
    font-size: 11px;
  }
  
  .chat-content {
    padding: 12px 4px;
    min-height: 0;
  }
  
  .message-container {
    padding: 0 4px;
  }
  
  .message-item {
    max-width: 90%;
  }
  
  .message-bubble {
    padding: 12px 16px;
    font-size: 14px;
    min-width: 30px; /* 移动端最小宽度稍小 */
  }
  
  .file-message {
    min-width: 180px; /* 移动端文件消息最小宽度 */
  }
  
  .chat-input-area {
    padding: 12px 16px;
    flex-shrink: 0;
  }
  
  .input-container {
    gap: 8px;
  }
  
  .send-btn {
    height: 32px;
    padding: 0 16px;
    font-size: 13px;
  }
}

/* 动画效果 */
.message-wrapper {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: white) {
  .chat-room-wrapper {
    background-color: #2c2c2c;
  }
  
  .chat-header {
    background: #3a3a3a;
    border-bottom-color: #484848;
  }
  
  .room-title {
    color: #ffffff;
  }
  
  .room-subtitle {
    color: #b3b3b3;
  }
  
  .chat-content {
    background: linear-gradient(to bottom, #2c2c2c 0%, #262626 100%);
  }
  
  .other-bubble {
    background: #404040;
    color: #ffffff;
  }
  
  .other-bubble::before {
    border-right-color: #404040;
  }
  
  .chat-input-area {
    background: #3a3a3a;
    border-top-color: #484848;
  }
  
  .message-input :deep(.el-textarea__inner) {
    background: #404040;
    border-color: #484848;
    color: #ffffff;
  }
  
  .toolbar-btn {
    background: #404040;
    color: #b3b3b3;
  }
  
  .toolbar-btn:hover {
    background: #484848;
    color: #409eff;
  }
}
</style>