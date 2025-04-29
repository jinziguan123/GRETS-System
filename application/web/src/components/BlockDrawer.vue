<template>
  <div class="block-drawer">
    <el-drawer
      v-model="visible"
      size="50%"
      title="区块链信息查询"
      direction="rtl"
      :destroy-on-close="false"
    >
      <div class="block-drawer-container">

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
                      <span class="label">交易数量:</span>
                      <span class="value">{{ block.transactionCount }}</span>
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

<script lang="ts" >
// import { defineComponent, ref, reactive, onMounted } from 'vue'
import { queryBlockList } from '@/api/block'
import { ElMessage } from 'element-plus'
import { DataAnalysis } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

interface BlockInfo {
  blockNumber: string
  blockHash: string
  prevHash: string
  dataHash: string
  saveTime: number
  transactionCount: number
}

export default defineComponent({
  name: 'BlockDrawer',
  components: {
    DataAnalysis
  },
  setup() {
    const visible = ref(false)
    const loading = ref(false)
    const blockList = ref<BlockInfo[]>([])
    
    const formState = reactive({
      blockNumber: ''
    })
    
    const pagination = reactive({
      current: 1,
      pageSize: 5,
      total: 0
    })
    
    const userStore = useUserStore()
    
    const showDrawer = () => {
      visible.value = true
      fetchBlockList()
    }
    
    const closeDrawer = () => {
      visible.value = false
    }
    
    const fetchBlockList = async () => {
      loading.value = true
      try {
        const params = {
          pageSize: pagination.pageSize,
          pageNumber: pagination.current,
          blockNumber: formState.blockNumber || undefined,
          organization: userStore.organization
        }
        
        const response = await queryBlockList(params)
        console.log(response)
        blockList.value = response.blocks
        pagination.total = response.total
      } catch (error) {
        console.error('获取区块列表失败:', error)
        ElMessage.error('获取区块列表失败')
        blockList.value = []
        pagination.total = 0
      } finally {
        loading.value = false
      }
    }
    
    const handleSubmit = () => {
      pagination.current = 1
      fetchBlockList()
    }
    
    const resetForm = () => {
      formState.blockNumber = ''
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
    
    const formatDate = (timestamp: number) => {
      if (!timestamp) return '-'
      const date = new Date(timestamp)
      return date.toLocaleString()
    }
    
    onMounted(() => {
      // 初始加载时不自动打开抽屉，但当用户点击悬浮按钮时会触发加载
    })
    
    return {
      visible,
      loading,
      blockList,
      formState,
      pagination,
      showDrawer,
      closeDrawer,
      fetchBlockList,
      handleSubmit,
      resetForm,
      handlePageChange,
      handleSizeChange,
      formatDate
    }
  }
})
</script>

<style scoped>
.block-drawer-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
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
</style> 