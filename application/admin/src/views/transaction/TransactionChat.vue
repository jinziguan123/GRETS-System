<template>
  <div class="chat-container">
    <el-row :gutter="20">
      <!-- 左侧交易列表 -->
      <el-col :span="6">
        <el-card class="chat-list-card">
          <template #header>
            <div class="card-header">
              <span>我的交易</span>
              <el-button type="primary" size="small" @click="refreshTransactions">刷新</el-button>
            </div>
          </template>
          
          <div v-if="loadingTransactions" class="loading-container">
            <el-skeleton :rows="3" animated />
          </div>
          
          <div v-else-if="transactions.length === 0" class="empty-data">
            <el-empty description="暂无交易" />
          </div>
          
          <div v-else class="transaction-list">
            <div 
              v-for="item in transactions" 
              :key="item.id"
              class="transaction-item"
              :class="{'active': currentTransaction && currentTransaction.id === item.id}"
              @click="selectTransaction(item)"
            >
              <div class="transaction-info">
                <div class="transaction-title">
                  {{ item.realtyTitle }}
                </div>
                <div class="transaction-meta">
                  <el-tag 
                    size="small" 
                    :type="getStatusType(item.status)"
                  >
                    {{ getStatusText(item.status) }}
                  </el-tag>
                  <span class="transaction-date">{{ formatDate(item.lastMessageTime) }}</span>
                </div>
              </div>
              <div v-if="item.unreadCount" class="unread-badge">{{ item.unreadCount }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 右侧聊天区域 -->
      <el-col :span="18">
        <el-card v-if="!currentTransaction" class="empty-chat-card">
          <div class="select-chat-prompt">
            <el-icon><ChatDotSquare /></el-icon>
            <span>请选择一个交易开始聊天</span>
          </div>
        </el-card>
        
        <el-card v-else class="chat-card">
          <!-- 聊天头部 -->
          <template #header>
            <div class="chat-header">
              <div class="chat-title">
                <h3>{{ currentTransaction.realtyTitle }}</h3>
                <el-tag 
                  size="small" 
                  :type="getStatusType(currentTransaction.status)"
                >
                  {{ getStatusText(currentTransaction.status) }}
                </el-tag>
              </div>
              <div class="chat-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="viewTransactionDetail"
                >
                  查看交易详情
                </el-button>
              </div>
            </div>
          </template>
          
          <!-- 聊天消息区域 -->
          <div class="chat-messages" ref="messagesContainer">
            <div v-if="loadingMessages" class="loading-container">
              <el-skeleton :rows="5" animated />
            </div>
            
            <div v-else-if="messages.length === 0" class="empty-data">
              <el-empty description="暂无消息，发送第一条消息开始商谈" />
            </div>
            
            <template v-else>
              <div v-for="(msg, index) in messages" :key="index" class="message-item">
                <!-- 日期分隔线 -->
                <div 
                  v-if="shouldShowDateDivider(msg, messages[index - 1])" 
                  class="date-divider"
                >
                  {{ formatDateDivider(msg.timestamp) }}
                </div>
                
                <!-- 消息内容 -->
                <div 
                  class="message" 
                  :class="{'message-mine': msg.fromId === userId, 'message-other': msg.fromId !== userId}"
                >
                  <div class="message-sender" v-if="msg.fromId !== userId">
                    {{ msg.anonymous ? '匿名用户' : msg.sender }}
                  </div>
                  <div class="message-content">
                    {{ msg.content }}
                    <span class="message-time">{{ formatMessageTime(msg.timestamp) }}</span>
                  </div>
                </div>
              </div>
            </template>
          </div>
          
          <!-- 聊天输入区域 -->
          <div class="chat-input">
            <el-input
              v-model="messageInput"
              type="textarea"
              :rows="3"
              placeholder="输入消息..."
              @keyup.enter.ctrl="sendMessage"
            />
            <div class="input-actions">
              <span class="input-tip">按 Ctrl+Enter 发送</span>
              <el-checkbox v-model="anonymous">匿名发送</el-checkbox>
              <el-button type="primary" @click="sendMessage" :disabled="!messageInput.trim()">
                发送
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { ChatDotSquare } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()
const userId = computed(() => userStore.userId)

// 交易列表
const transactions = ref([])
const loadingTransactions = ref(false)
const currentTransaction = ref(null)

// 聊天消息
const messages = ref([])
const loadingMessages = ref(false)
const messageInput = ref('')
const messagesContainer = ref(null)
const anonymous = ref(false)

// 获取交易列表
const refreshTransactions = async () => {
  loadingTransactions.value = true
  try {
    // 这里替换为实际的API调用
    setTimeout(() => {
      // 模拟数据
      transactions.value = [
        {
          id: 'T001',
          realtyId: 'R001',
          realtyTitle: '浦东新区 | 阳光花园三室两厅',
          status: 'negotiating',
          buyerId: 'B001',
          sellerId: 'S001',
          lastMessageTime: new Date(Date.now() - 30 * 60 * 1000),
          unreadCount: 3
        },
        {
          id: 'T002',
          realtyId: 'R002',
          realtyTitle: '长宁区 | 新华花苑两室一厅',
          status: 'pending_payment',
          buyerId: 'B002',
          sellerId: 'S002',
          lastMessageTime: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
          unreadCount: 0
        },
        {
          id: 'T003',
          realtyId: 'R003',
          realtyTitle: '静安区 | 静安豪庭四室两厅',
          status: 'completed',
          buyerId: 'B001',
          sellerId: 'S003',
          lastMessageTime: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
          unreadCount: 0
        }
      ]
      loadingTransactions.value = false
    }, 800)
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
    ElMessage.error('获取交易列表失败')
    loadingTransactions.value = false
  }
}

// 加载聊天消息
const loadMessages = async (transactionId) => {
  loadingMessages.value = true
  messages.value = []
  
  try {
    // 这里替换为实际的API调用
    setTimeout(() => {
      // 模拟数据
      const now = new Date()
      messages.value = [
        {
          id: 'm1',
          transactionId,
          fromId: 'B001',
          sender: '买家',
          content: '您好，我对这套房子很感兴趣，价格能商量吗？',
          timestamp: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000 - 30 * 60 * 1000),
          anonymous: true
        },
        {
          id: 'm2',
          transactionId,
          fromId: 'S001',
          sender: '卖家',
          content: '您好，可以谈，您预期多少？',
          timestamp: new Date(now.getTime() - 3 * 24 * 60 * 60 * 1000 - 25 * 60 * 1000),
          anonymous: true
        },
        {
          id: 'm3',
          transactionId,
          fromId: 'B001',
          sender: '买家',
          content: '我看市场价大概在500万左右，我希望480万成交，可以吗？',
          timestamp: new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000 - 15 * 60 * 1000),
          anonymous: true
        },
        {
          id: 'm4',
          transactionId,
          fromId: 'S001',
          sender: '卖家',
          content: '这个价格有点低，考虑一下495万？',
          timestamp: new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000 - 5 * 60 * 1000),
          anonymous: true
        },
        {
          id: 'm5',
          transactionId,
          fromId: 'B001',
          sender: '买家',
          content: '我可以考虑490万，这是我的最高价位了。',
          timestamp: new Date(now.getTime() - 50 * 60 * 1000),
          anonymous: false
        },
        {
          id: 'm6',
          transactionId,
          fromId: 'S001',
          sender: '卖家',
          content: '成交。我们可以准备合同了。',
          timestamp: new Date(now.getTime() - 30 * 60 * 1000),
          anonymous: false
        }
      ]
      
      loadingMessages.value = false
      // 滚动到最底部
      nextTick(() => {
        scrollToBottom()
      })
    }, 1000)
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('获取聊天记录失败')
    loadingMessages.value = false
  }
}

// 选择交易
const selectTransaction = (transaction) => {
  currentTransaction.value = transaction
  loadMessages(transaction.id)
}

// 发送消息
const sendMessage = () => {
  if (!messageInput.value.trim()) return
  
  if (!currentTransaction.value) {
    ElMessage.warning('请先选择一个交易')
    return
  }
  
  // 这里替换为实际的API调用
  const newMessage = {
    id: `m${Date.now()}`,
    transactionId: currentTransaction.value.id,
    fromId: userId.value,
    sender: userStore.username,
    content: messageInput.value,
    timestamp: new Date(),
    anonymous: anonymous.value
  }
  
  messages.value.push(newMessage)
  messageInput.value = ''
  
  // 滚动到最底部
  nextTick(() => {
    scrollToBottom()
  })
}

// 查看交易详情
const viewTransactionDetail = () => {
  if (currentTransaction.value) {
    router.push(`/transaction/${currentTransaction.value.id}`)
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 获取交易状态类型
const getStatusType = (status) => {
  const statusMap = {
    'negotiating': 'warning',
    'pending_payment': 'primary',
    'processing': 'success',
    'completed': 'info',
    'canceled': 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取交易状态文本
const getStatusText = (status) => {
  const statusMap = {
    'negotiating': '商谈中',
    'pending_payment': '待支付',
    'processing': '处理中',
    'completed': '已完成',
    'canceled': '已取消'
  }
  return statusMap[status] || '未知状态'
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  const now = new Date()
  const messageDate = new Date(date)
  
  if (now.toDateString() === messageDate.toDateString()) {
    return dayjs(date).format('HH:mm')
  } else if (now.getFullYear() === messageDate.getFullYear()) {
    return dayjs(date).format('MM-DD')
  } else {
    return dayjs(date).format('YYYY-MM-DD')
  }
}

// 是否显示日期分隔线
const shouldShowDateDivider = (currentMsg, prevMsg) => {
  if (!prevMsg) return true
  
  const currentDate = new Date(currentMsg.timestamp)
  const prevDate = new Date(prevMsg.timestamp)
  
  return currentDate.toDateString() !== prevDate.toDateString()
}

// 格式化日期分隔线
const formatDateDivider = (date) => {
  if (!date) return ''
  
  const messageDate = new Date(date)
  const now = new Date()
  const yesterday = new Date(now)
  yesterday.setDate(now.getDate() - 1)
  
  if (now.toDateString() === messageDate.toDateString()) {
    return '今天'
  } else if (yesterday.toDateString() === messageDate.toDateString()) {
    return '昨天'
  } else {
    return dayjs(date).format('YYYY年MM月DD日')
  }
}

// 格式化消息时间
const formatMessageTime = (date) => {
  if (!date) return ''
  return dayjs(date).format('HH:mm')
}

// 初始化
onMounted(() => {
  refreshTransactions()
})
</script>

<style scoped>
.chat-container {
  padding: 20px;
}

.loading-container {
  padding: 20px 0;
}

.empty-data {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 0;
}

.chat-list-card {
  height: calc(100vh - 160px);
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.transaction-list {
  flex: 1;
  overflow-y: auto;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background-color 0.3s;
}

.transaction-item:hover, .transaction-item.active {
  background-color: #f5f7fa;
}

.transaction-info {
  flex: 1;
  overflow: hidden;
}

.transaction-title {
  margin-bottom: 8px;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.transaction-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.transaction-date {
  margin-left: 8px;
}

.unread-badge {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 22px;
  height: 22px;
  background-color: #f56c6c;
  color: white;
  border-radius: 50%;
  font-size: 12px;
}

.empty-chat-card {
  height: calc(100vh - 160px);
  display: flex;
  justify-content: center;
  align-items: center;
}

.select-chat-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #909399;
}

.select-chat-prompt .el-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.chat-card {
  height: calc(100vh - 160px);
  display: flex;
  flex-direction: column;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-title {
  display: flex;
  align-items: center;
}

.chat-title h3 {
  margin: 0;
  margin-right: 10px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
}

.date-divider {
  text-align: center;
  margin: 10px 0;
  position: relative;
}

.date-divider::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  width: 100%;
  height: 1px;
  background-color: #ebeef5;
  z-index: 1;
}

.date-divider span {
  position: relative;
  background-color: white;
  padding: 0 10px;
  font-size: 12px;
  color: #909399;
  z-index: 2;
}

.message-item {
  margin-bottom: 15px;
}

.message {
  display: flex;
  flex-direction: column;
  max-width: 70%;
}

.message-mine {
  align-self: flex-end;
  margin-left: auto;
}

.message-other {
  align-self: flex-start;
}

.message-sender {
  font-size: 12px;
  color: #909399;
  margin-bottom: 5px;
}

.message-content {
  position: relative;
  padding: 10px 15px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.message-mine .message-content {
  background-color: #ecf5ff;
  color: #409eff;
  border-top-right-radius: 0;
}

.message-other .message-content {
  background-color: #f5f5f5;
  color: #333;
  border-top-left-radius: 0;
}

.message-time {
  position: absolute;
  bottom: -18px;
  right: 5px;
  font-size: 10px;
  color: #909399;
}

.message-mine .message-time {
  right: 5px;
}

.message-other .message-time {
  right: 5px;
}

.chat-input {
  margin-top: auto;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.input-tip {
  font-size: 12px;
  color: #909399;
}
</style>