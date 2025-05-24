<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <el-aside width="220px" class="sidebar">
      <div class="logo-container">
        <h1 class="logo">GRETS后台管理</h1>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        class="el-menu-vertical"
        :collapse="isCollapse"
        :router="true"
        background-color="#001529"
        text-color="#fff"
        active-text-color="#409EFF"
      >
        <el-menu-item index="/">
          <el-icon><Monitor /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        
        <!-- 房产管理 - 所有组织可见 -->
        <el-sub-menu index="/realty">
          <template #title>
            <el-icon><House /></el-icon>
            <span>房产管理</span>
          </template>
          <el-menu-item index="/realty">房产列表</el-menu-item>
          <!-- 政府组织可创建房产 -->
          <el-menu-item index="/realty/create" v-if="hasOrganization('government')">添加房产</el-menu-item>
        </el-sub-menu>
        
        <!-- 交易管理 - 投资者、银行、政府可见 -->
        <el-sub-menu index="/transaction" v-if="hasOrganization(['investor', 'bank', 'government', 'audit'])">
          <template #title>
            <el-icon><Sell /></el-icon>
            <span>交易管理</span>
          </template>
          <el-menu-item index="/transaction">交易列表</el-menu-item>
<!--          <el-menu-item index="/transaction/create" v-if="hasOrganization('investor')">创建交易</el-menu-item>-->
          <!-- 投资者可见 - 聊天室 -->
<!--          <el-menu-item index="/transaction/chat" v-if="hasOrganization('investor')">交易聊天室</el-menu-item>-->
        </el-sub-menu>
        
        <!-- 合同管理 - 投资者、银行、政府、审计可见 -->
        <el-sub-menu index="/contract" v-if="hasOrganization(['investor', 'bank', 'government', 'audit'])">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>合同管理</span>
          </template>
          <el-menu-item index="/contract">合同列表</el-menu-item>
          <el-menu-item index="/contract/create" v-if="hasOrganization(['investor', 'government'])">创建合同</el-menu-item>
          <!-- 审计组织特有 - 合规审核 -->
<!--          <el-menu-item index="/contract/audit" v-if="hasOrganization('audit')">合规审核</el-menu-item>-->
        </el-sub-menu>
        
        <!-- 支付管理 - 投资者、银行可见 -->
        <el-sub-menu index="/payment" v-if="hasOrganization(['investor', 'bank'])">
          <template #title>
            <el-icon><Money /></el-icon>
            <span>支付管理</span>
          </template>
          <el-menu-item index="/payment">支付列表</el-menu-item>
        </el-sub-menu>

        <!-- 聊天室管理 - 投资者可见 -->
        <el-sub-menu index="/chat" v-if="hasOrganization(['investor', 'government'])">
          <template #title>
            <el-icon><CreditCard /></el-icon>
            <span>聊天室管理</span>
          </template>
          <el-menu-item index="/chat">聊天室列表</el-menu-item>
        </el-sub-menu>

<!--        &lt;!&ndash; 税费管理 - 政府和投资者可见 &ndash;&gt;-->
<!--        <el-sub-menu index="/tax" v-if="hasOrganization(['government', 'investor'])">-->
<!--          <template #title>-->
<!--            <el-icon><Wallet /></el-icon>-->
<!--            <span>税费管理</span>-->
<!--          </template>-->
<!--          <el-menu-item index="/tax">税费列表</el-menu-item>-->
<!--          <el-menu-item index="/tax/create" v-if="hasOrganization('government')">创建税费</el-menu-item>-->
<!--        </el-sub-menu>-->
        
        <!-- 抵押贷款 - 投资者和银行可见 -->
<!--        <el-sub-menu index="/mortgage" v-if="hasOrganization(['investor', 'bank'])">-->
<!--          <template #title>-->
<!--            <el-icon><CreditCard /></el-icon>-->
<!--            <span>抵押贷款</span>-->
<!--          </template>-->
<!--          <el-menu-item index="/mortgage">贷款列表</el-menu-item>-->
<!--          <el-menu-item index="/mortgage/create" v-if="hasOrganization('investor')">申请贷款</el-menu-item>-->
<!--          <el-menu-item index="/mortgage/approve" v-if="hasOrganization('bank')">贷款审批</el-menu-item>-->
<!--        </el-sub-menu>-->
        
        <!-- 统计分析 - 政府和银行可见 -->
        <el-sub-menu index="/statistics" v-if="hasOrganization(['government'])">
          <template #title>
            <el-icon><DataAnalysis /></el-icon>
            <span>统计分析</span>
          </template>
          <el-menu-item index="/statistics/transaction" v-if="hasOrganization(['government'])">交易统计</el-menu-item>
<!--          <el-menu-item index="/statistics/loan" v-if="hasOrganization('bank')">贷款统计</el-menu-item>-->
        </el-sub-menu>
        
        <!-- 系统管理 - 只有管理员可见 -->
        <el-sub-menu index="/admin" v-if="hasOrganization('administrator')">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin/users">用户管理</el-menu-item>
          <el-menu-item index="/admin/system">系统设置</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- 区块链信息查询 -->
    <BlockDrawer />
    
    <!-- 主要内容区域 -->
    <el-container class="main-container">
      <el-header class="header">
        <div class="header-left">
          <el-button type="text" @click="toggleSidebar">
            <el-icon><Menu /></el-icon>
          </el-button>
        </div>
        <div class="header-right">
          <div class="user-info">
            <el-tag
                class="organization-tag"
                :type="getOrganizationTagType(userOrganization)"
                size="small"
            >
              {{ organizationName }}
            </el-tag>
            <span class="username" style="margin-left: 5px">{{ username }}</span>
          </div>
          <el-dropdown trigger="click">
            <span class="user-dropdown">
              <el-badge :is-dot="hasPendingTransactions" class="notification-badge">
                <el-avatar :size="32" :icon="UserFilled"></el-avatar>
              </el-badge>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-badge :is-dot="hasPendingTransactions" class="notification-badge">
                  <el-dropdown-item @click="router.push('/user/profile')">个人信息</el-dropdown-item>
                </el-badge>
<!--                <el-dropdown-item @click="router.push('/change-password')">修改密码</el-dropdown-item>-->
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <!-- 使用MainContent组件替代原有内容 -->
        <main-content />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
// import { ref, computed } from 'vue'
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import MainContent from '@/components/MainContent.vue'
import BlockDrawer from '@/components/BlockDrawer.vue'
import {
  Menu, 
  House, 
  Document, 
  Sell, 
  Money,
  CreditCard, 
  Setting, 
  Monitor,
  ArrowDown,
  DataAnalysis,
  UserFilled,
} from '@element-plus/icons-vue'
import { queryTransactionList } from '@/api/transaction'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 用户信息
const username = computed(() => userStore.user?.name || '')
const userOrganization = computed(() => userStore.user?.organization || '')

// 待处理交易状态（有则显示小红点）
const hasPendingTransactions = ref(false)

// 获取待处理交易数据
const fetchPendingTransactions = async () => {
  try {
    // 这里应该调用相应的API获取待处理交易
    const response = await queryTransactionList({
      status: 'PENDING',
      sellerCitizenID: userStore.user?.citizenID,
      sellerOrganization: userStore.user?.organization,
      pageSize: 1000,
      pageNumber: 1
    })
    hasPendingTransactions.value = response.transactions.length > 0
  } catch (error) {
    console.error('获取待处理交易失败:', error)
  }
}

// 组件挂载时获取待处理交易
onMounted(() => {
  fetchPendingTransactions()
})

// 组织名称和样式
const organizationName = computed(() => {
  const orgMap: Record<string, string> = {
    'administrator': '系统管理员',
    'government': '政府监管部门',
    'investor': '投资者/买家',
    'bank': '银行机构',
    'audit': '审计监管部门',
    'thirdparty': '第三方机构'
  }
  return orgMap[userOrganization.value] || '未知组织'
})

// 侧边栏折叠状态
const isCollapse = ref<boolean>(false)
const toggleSidebar = (): void => {
  isCollapse.value = !isCollapse.value
}

// 当前激活的菜单项
const activeMenu = computed(() => {
  return route.path
})

// 组织权限检查
const hasOrganization = (orgs: string | string[]): boolean => {
  if (!orgs) return true
  if (typeof orgs === 'string') {
    return userOrganization.value === orgs
  }
  return orgs.includes(userOrganization.value)
}

// 退出登录
const logout = (): void => {
  userStore.logout()
  router.push('/auth/login')
}

// 组织样式
const getOrganizationTagType = (org: string): 'success' | 'warning' | 'danger' | 'info' | 'primary' | undefined => {
  const typeMap: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary' | undefined> = {
    'administrator': 'primary',
    'government': 'danger',
    'investor': 'success',
    'bank': 'warning',
    'audit': 'info',
    'thirdparty': 'primary'
  }
  return typeMap[org]
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
  display: flex;
}

.sidebar {
  height: 100%;
  background-color: #001529;
  transition: width 0.3s;
  overflow-x: hidden;
}

.logo-container {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0 20px;
  color: #fff;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  font-size: 20px;
  font-weight: bold;
  margin: 0;
}

.el-menu-vertical {
  border-right: none;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.header-left, .header-right {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-right: 10px;
}

.organization-tag {
  font-weight: normal;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 8px;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f7fa;
}

.organization-badge {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  margin-right: 15px;
  color: white;
}

.org-administrator {
  background-color: #409EFF;
}

.org-government {
  background-color: #67C23A;
}

.org-investor {
  background-color: #E6A23C;
}

.org-bank {
  background-color: #F56C6C;
}

.org-audit {
  background-color: #909399;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    z-index: 999;
    height: 100%;
  }
  
  .main-container {
    margin-left: 0;
  }
}

/* 消息通知小红点样式 */
.notification-badge :deep(.el-badge__content.is-dot) {
  right: 5px;
  top: 5px;
}
</style> 