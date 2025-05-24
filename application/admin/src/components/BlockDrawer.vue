<template>
  <div class="block-drawer">
    <el-drawer
      v-model="visible"
      :size="!detailVisible ? '50%' : '100%'"
      title="区块链信息查询"
      direction="rtl"
      :destroy-on-close="false"
      :modal="!detailVisible"
      :append-to-body="true"
      :class="{ 'drawer-inactive': detailVisible }"
      :with-header="true"
      :show-close="!detailVisible"
      :close-on-click-modal="!detailVisible"
      :close-on-press-escape="!detailVisible"
    >
      <div class="block-drawer-container">
        <div class="search-form">
          <el-form :model="formState" inline @submit.prevent="handleSubmit">
            <el-form-item label="区块哈希">
              <el-input 
                v-model="formState.blockHash" 
                placeholder="输入区块哈希" 
                clearable
                :disabled="detailVisible"
              />
            </el-form-item>
            <el-form-item label="所属省份">
              <el-select 
                v-model="formState.provinceName" 
                placeholder="选择省份" 
                clearable
                style="width: 200px"
                :disabled="detailVisible"
              >
                <el-option 
                  v-for="province in provinceList" 
                  :key="province.provinceCode" 
                  :label="province.provinceName" 
                  :value="province.provinceName" 
                />
              </el-select>
            </el-form-item>
            <el-form-item label="创建者">
              <el-input 
                v-model="formState.creator" 
                placeholder="输入交易创建者" 
                clearable
                :disabled="detailVisible"
                style="width: 200px"
              />
            </el-form-item>
            <div>
              <el-button type="primary" @click="handleSubmit" :disabled="detailVisible">查询</el-button>
              <el-button @click="resetForm" :disabled="detailVisible">重置</el-button>
            </div>
          </el-form>
        </div>

        <div class="block-list-container" v-loading="loading" element-loading-text="加载中...">
          <el-empty v-if="blockList.length === 0 && !loading" description="暂无数据" />
          <template v-else>
            <div class="block-item" v-for="block in blockList" :key="block.blockNumber">
              <el-collapse>
                <el-collapse-item :name="block.blockNumber">
                  <template #title>
                    <div class="block-header">
                      <span>区块号: {{ block.blockNumber }}</span>
                      <span>时间: {{ formatDate(block.saveTime) }}</span>
                    </div>
                  </template>
                  <div class="block-details">
                    <div class="detail-item">
                      <span class="label">通道名称:</span>
                      <span class="value">{{ block.channelName || '未知' }}</span>
                    </div>
                    <div class="detail-item">
                      <span class="label">交易数量:</span>
                      <span class="value">{{ block.txCount }}</span>
                    </div>
                    <div class="detail-item">
                      <span class="label">区块哈希:</span>
                      <span class="value">{{ block.blockHash }}</span>
                    </div>
                    <div class="detail-item">
                      <span class="label">上一区块哈希:</span>
                      <span class="value">{{ block.prevHash }}</span>
                    </div>
                    <div class="detail-item">
                      <span class="label">创建时间:</span>
                      <span class="value">{{ formatDate(block.saveTime) }}</span>
                    </div>
                    <div class="detail-item">
                      <span class="label">数据哈希:</span>
                      <span class="value">{{ block.dataHash }}</span>
                    </div>
                    <div class="action-item">
                      <el-button type="primary" size="small" @click.stop="showBlockDetail(block)">
                        查看区块交易
                      </el-button>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
            <div class="pagination-container">
              <el-pagination
                v-model:current-page="pagination.current"
                :page-size="pagination.pageSize"
                :total="pagination.total"
                @current-change="handlePageChange"
                layout="total, sizes, prev, pager, next"
                :page-sizes="[5, 10, 20]"
                @size-change="handleSizeChange"
              />
            </div>
          </template>
        </div>
      </div>
    </el-drawer>

    <!-- 区块交易详情抽屉 -->
    <el-drawer
      v-model="detailVisible"
      size="50%"
      title="区块交易详情"
      direction="rtl"
      :destroy-on-close="false"
      :append-to-body="true"
      :before-close="handleDetailClose"
    >
      <div class="block-detail-container" v-loading="detailLoading" element-loading-text="加载中...">
        <div v-if="currentBlock" class="block-info">
          <h3>区块基本信息</h3>
          <div class="info-container">
            <div class="info-item">
              <span class="label">区块号:</span>
              <span class="value">{{ currentBlock.blockNumber }}</span>
            </div>
            <div class="info-item">
              <span class="label">通道名称:</span>
              <span class="value">{{ currentBlock.channelName || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="label">区块哈希:</span>
              <span class="value">{{ currentBlock.blockHash }}</span>
            </div>
            <div class="info-item">
              <span class="label">上一区块哈希:</span>
              <span class="value">{{ currentBlock.prevHash }}</span>
            </div>
            <div class="info-item">
              <span class="label">数据哈希:</span>
              <span class="value">{{ currentBlock.dataHash }}</span>
            </div>
            <div class="info-item">
              <span class="label">创建时间:</span>
              <span class="value">{{ formatDate(currentBlock.saveTime) }}</span>
            </div>
          </div>
        </div>

        <div class="transaction-list">
          <h3>交易列表</h3>
          <el-empty v-if="transactions.length === 0 && !detailLoading" description="暂无交易数据" />
          <template v-else>
            <div class="transaction-item" v-for="(tx, index) in transactions" :key="index">
              <el-collapse>
                <el-collapse-item :title="`交易ID: ${tx.transactionID}`" style="margin-left: 10px">
                  <div class="tx-details">
                    <div class="tx-item">
                      <span class="label">交易ID:</span>
                      <span class="value">{{ tx.transactionID }}</span>
                    </div>
                    <div class="tx-item">
                      <span class="label">时间戳:</span>
                      <span class="value">{{ formatDate(tx.transactionTimestamp) }}</span>
                    </div>
                    <div class="tx-item">
                      <span class="label">创建者:</span>
                      <span class="value">{{ tx.creator || '-' }}</span>
                    </div>
                    <div class="tx-item">
                      <span class="label">函数名称:</span>
                      <span class="value">{{ tx.chainCodeFunctionName || '-' }}</span>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
          </template>
        </div>
      </div>
    </el-drawer>

    <div class="float-button-container" v-if="!visible">
      <el-popover
          class="box-item"
          title="区块详情"
          content="点击查看区块信息"
          placement="left"
      >
        <template #reference>
          <el-button
              type="primary"
              circle
              class="float-button"
              @click="showDrawer"
              size="large"
          >
            <el-icon><DataAnalysis /></el-icon>
          </el-button>
        </template>
      </el-popover>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue'
import { queryBlockList, queryBlockTransactionList } from '@/api/block'
import { ElMessage } from 'element-plus'
import { DataAnalysis } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

interface BlockInfo {
  blockNumber: string
  blockHash: string
  prevHash: string
  dataHash: string
  saveTime: number
  txCount: number
  channelName?: string
}

interface Transaction {
  transactionID: string
  creator: string
  transactionTimestamp: string
  chainCodeFunctionName: string
}

interface Province {
  provinceCode: string
  provinceName: string
}

const visible = ref(false)
const detailVisible = ref(false)
const loading = ref(false)
const detailLoading = ref(false)
const blockList = ref<BlockInfo[]>([])
const transactions = ref<Transaction[]>([])
const currentBlock = ref<BlockInfo | null>(null)

const formState = reactive({
  blockHash: '',
  provinceName: '',
  creator: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 5,
  total: 0
})

const userStore = useUserStore()

const provinceList = ref<Province[]>([
  { provinceCode: '11', provinceName: '北京市' },
  { provinceCode: '12', provinceName: '天津市' },
  { provinceCode: '13', provinceName: '河北省' },
  { provinceCode: '14', provinceName: '山西省' },
  { provinceCode: '15', provinceName: '内蒙古自治区' },
  { provinceCode: '21', provinceName: '辽宁省' },
  { provinceCode: '22', provinceName: '吉林省' },
  { provinceCode: '23', provinceName: '黑龙江省' },
  { provinceCode: '31', provinceName: '上海市' },
  { provinceCode: '32', provinceName: '江苏省' },
  { provinceCode: '33', provinceName: '浙江省' },
  { provinceCode: '34', provinceName: '安徽省' },
  { provinceCode: '35', provinceName: '福建省' },
  { provinceCode: '36', provinceName: '江西省' },
  { provinceCode: '37', provinceName: '山东省' },
  { provinceCode: '41', provinceName: '河南省' },
  { provinceCode: '42', provinceName: '湖北省' },
  { provinceCode: '43', provinceName: '湖南省' },
  { provinceCode: '44', provinceName: '广东省' },
  { provinceCode: '45', provinceName: '广西壮族自治区' },
  { provinceCode: '46', provinceName: '海南省' },
  { provinceCode: '50', provinceName: '重庆市' },
  { provinceCode: '51', provinceName: '四川省' },
  { provinceCode: '52', provinceName: '贵州省' },
  { provinceCode: '53', provinceName: '云南省' },
  { provinceCode: '54', provinceName: '西藏自治区' },
  { provinceCode: '61', provinceName: '陕西省' },
  { provinceCode: '62', provinceName: '甘肃省' },
  { provinceCode: '63', provinceName: '青海省' },
  { provinceCode: '64', provinceName: '宁夏回族自治区' },
  { provinceCode: '65', provinceName: '新疆维吾尔自治区' },
  { provinceCode: '71', provinceName: '台湾省' },
  { provinceCode: '81', provinceName: '香港特别行政区' },
  { provinceCode: '82', provinceName: '澳门特别行政区' }
])

const showDrawer = () => {
  visible.value = true
  fetchBlockList()
}

const closeDrawer = () => {
  visible.value = false
}

const showBlockDetail = async (block: BlockInfo) => {
  currentBlock.value = block
  detailVisible.value = true
  await fetchBlockTransactions(block)
}

const handleDetailClose = (done: () => void) => {
  detailVisible.value = false
  transactions.value = []
  done()
}

const fetchBlockList = async () => {
  loading.value = true
  try {
    const params = {
      pageSize: pagination.pageSize,
      pageNumber: pagination.current,
      creator: formState.creator || undefined,
      blockHash: formState.blockHash || undefined,
      provinceName: formState.provinceName || undefined,
      organization: userStore.organization
    }
    
    const response = await queryBlockList(params)
    blockList.value = response.blocks || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取区块列表失败:', error)
    blockList.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

const fetchBlockTransactions = async (block: BlockInfo) => {
  if (!block.blockNumber || !block.channelName) {
    ElMessage.warning('区块信息不完整，无法获取交易详情')
    return
  }
  
  detailLoading.value = true
  transactions.value = []
  
  try {
    const params = {
      blockNumber: block.blockNumber,
      channelName: block.channelName,
      organization: userStore.organization
    }
    
    const response = await queryBlockTransactionList(params)
    transactions.value = response.transactionList || []
  } catch (error) {
    console.error('获取区块交易详情失败:', error)
    ElMessage.error('获取区块交易详情失败')
  } finally {
    detailLoading.value = false
  }
}

const handleSubmit = () => {
  pagination.current = 1
  fetchBlockList()
}

const resetForm = () => {
  formState.blockHash = ''
  formState.provinceName = ''
  formState.creator = ''
  pagination.current = 1
  fetchBlockList()
}

const handlePageChange = (page: number) => {
  pagination.current = page
  fetchBlockList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.current = 1
  fetchBlockList()
}

const formatDate = (timestamp: string) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString()
}

onMounted(() => {
  // 初始加载时不自动打开抽屉，但当用户点击悬浮按钮时会触发加载
})
</script>

<style scoped>
.block-drawer-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 16px;
}

.search-form {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 16px;
}

.block-list-container {
  margin-top: 16px;
}

.block-item {
  padding: 16px;
  margin-bottom: 12px;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid #f0f0f0;
}

.block-header {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.block-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  display: flex;
  line-height: 24px;
}

.action-item {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

.label {
  color: #8c8c8c;
  width: 100px;
  flex-shrink: 0;
}

.value {
  word-break: break-all;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.float-button-container {
  position: absolute;
  right: 24px;
  bottom: 50%;
  z-index: 10;
}

.float-button {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 区块详情样式 */
.block-detail-container {
  padding: 0 16px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.block-info {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
}

.info-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 16px;
}

.info-item {
  display: flex;
  line-height: 24px;
}

.transaction-list {
  margin-top: 16px;
}

.transaction-item {
  margin-bottom: 12px;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
}

.tx-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
}

.tx-item {
  display: flex;
  line-height: 24px;
}

.arg-item {
  padding: 4px 0;
}

/* 当详情抽屉打开时，原抽屉的灰色半透明效果 */
.drawer-inactive {
  pointer-events: none;
}

.drawer-inactive :deep(.el-drawer__body) {
  opacity: 0.5;
  filter: grayscale(50%);
}
</style> 