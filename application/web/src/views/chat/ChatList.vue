<template>
  <div class="chat-list-container">
    <div class="page-header">
      <h3>我的聊天室</h3>
      <div class="header-actions">
        <el-button @click="refreshList">刷新</el-button>
      </div>
    </div>

    <!-- 搜索筛选 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="房产证号">
          <el-input v-model="searchForm.realtyCert" placeholder="请输入房产证号" clearable />
        </el-form-item>
        <el-form-item label="状态" style="width: 200px;">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="进行中" value="ACTIVE" />
            <el-option label="已关闭" value="CLOSED" />
            <el-option label="已冻结" value="FROZEN" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 聊天室列表 -->
    <el-card class="list-card">
      <div v-loading="loading">
        <div v-if="chatRooms.length === 0 && !loading" class="empty-state">
          <el-empty description="暂无聊天室" />
        </div>
        <div v-else>
          <div 
            v-for="room in chatRooms" 
            :key="room.roomUUID"
            class="chat-room-item"
            :class="{ 'room-disabled': room.status !== 'ACTIVE' }"
            @click="enterChatRoom(room)"
          >
            <div class="room-avatar">
              <el-avatar :size="60">
                <el-icon><ChatDotRound /></el-icon>
              </el-avatar>
            </div>
            <div class="room-content">
              <div class="room-header">
                <div class="room-title">
                  <span class="title">{{ getRoomTitle(room) }}</span>
                  <el-tag :type="getRoomStatusType(room.status)" size="small">
                    {{ getRoomStatusText(room.status) }}
                  </el-tag>
                </div>
                <div class="room-time">
                  {{ formatTime(room.lastMessageTime) }}
                </div>
              </div>
              <div class="room-info">
                <div class="room-subtitle">
                  {{ getRoomSubtitle(room) }}
                </div>
                <div class="verification-amount">
                  验资金额: ¥{{ room.verificationAmount.toLocaleString() }}
                </div>
              </div>
              <div class="room-footer">
                <div class="last-message">
                  {{ ('最后消息: ' + room.lastMessageContent) || '暂无消息' }}
                </div>
                <div class="unread-count" v-if="getUnreadCount(room) > 0">
                  <el-badge :value="getUnreadCount(room)" :max="99" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          v-model:current-page="searchForm.pageNumber"
          v-model:page-size="searchForm.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ChatDotRound } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getChatRoomList } from '@/api/chat'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const loading = ref(false)
const chatRooms = ref([])
const total = ref(0)

// 搜索表单
const searchForm = reactive({
  realtyCert: '',
  status: '',
  pageSize: 10,
  pageNumber: 1
})

// 获取聊天室列表
const fetchChatRooms = async () => {
  loading.value = true
  try {
    const requestData = {
      userCitizenID: userStore.user.citizenID,
      userOrganization: userStore.user.organization,
      pageSize: searchForm.pageSize,
      pageNumber: searchForm.pageNumber
    }

    // 添加筛选条件
    if (searchForm.status) {
      requestData.status = searchForm.status
    }
    if (searchForm.realtyCert) {
      requestData.realtyCert = searchForm.realtyCert
    }

    const response = await getChatRoomList(requestData)
    chatRooms.value = response.chatRooms || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取聊天室列表失败:', error)
    ElMessage.error('获取聊天室列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  searchForm.pageNumber = 1
  fetchChatRooms()
}

// 重置搜索
const resetSearch = () => {
  searchForm.realtyCert = ''
  searchForm.status = ''
  searchForm.pageNumber = 1
  fetchChatRooms()
}

// 刷新列表
const refreshList = () => {
  fetchChatRooms()
}

// 改变每页显示数量
const handleSizeChange = (size) => {
  searchForm.pageSize = size
  fetchChatRooms()
}

// 改变页码
const handleCurrentChange = (page) => {
  searchForm.pageNumber = page
  fetchChatRooms()
}

// 进入聊天室
const enterChatRoom = (room) => {
  // 如果聊天室已关闭或冻结，显示提示信息
  if (room.status !== 'ACTIVE') {
    const statusText = room.status === 'CLOSED' ? '已关闭' : '已冻结'
    ElMessage.info(`聊天室${statusText}，您只能查看历史消息记录`)
  }
  
  router.push({
    path: `/chat/room/${room.roomUUID}`
  })
}

// 获取房间标题
const getRoomTitle = (room) => {
  if (room.realtyInfo) {
    return `${room.realtyInfo.address} - ${room.realtyInfo.realtyCert}`
  }
  return `房产 - ${room.realtyCert}`
}

// 获取房间副标题
const getRoomSubtitle = (room) => {
  if (room.realtyInfo) {
    return `${room.realtyInfo.address} · ${room.realtyInfo.area}m²`
  }
  return '房产交易咨询'
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

// 获取未读消息数
const getUnreadCount = (room) => {
  // 根据当前用户身份获取对应的未读数
  const currentUserHash = userStore.user.citizenIDHash
  if (room.buyerCitizenIDHash === currentUserHash) {
    return room.unreadCountBuyer
  } else if (room.sellerCitizenIDHash === currentUserHash) {
    return room.unreadCountSeller
  }
  return 0
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
    return date.toLocaleDateString()
  }
}

// 生命周期
onMounted(() => {
  fetchChatRooms()
})
</script>

<style scoped>
.chat-list-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
  font-size: 22px;
}

.search-card {
  margin-bottom: 20px;
}

.list-card {
  min-height: 400px;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

.chat-room-item {
  display: flex;
  padding: 16px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.chat-room-item:hover {
  background-color: #f8f9fa;
}

.chat-room-item:last-child {
  border-bottom: none;
}

.room-disabled {
  opacity: 0.7;
  background-color: #f8f9fa;
}

.room-disabled:hover {
  background-color: #f0f0f0;
}

.room-disabled .title {
  color: #909399;
}

.room-disabled .room-subtitle {
  color: #c0c4cc;
}

.room-disabled .last-message {
  color: #c0c4cc;
}

.room-avatar {
  margin-right: 16px;
}

.room-content {
  flex: 1;
  min-width: 0;
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.room-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.room-time {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
}

.room-info {
  margin-bottom: 8px;
}

.room-subtitle {
  font-size: 14px;
  color: #606266;
  margin-bottom: 4px;
}

.verification-amount {
  font-size: 12px;
  color: #409eff;
}

.room-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.last-message {
  font-size: 14px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: calc(100% - 60px);
}

.unread-count {
  margin-left: 8px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style> 