<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <el-aside width="220px" class="sidebar">
      <div class="logo-container">
        <h1 class="logo">GRETS</h1>
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
        
        <el-sub-menu index="/realty">
          <template #title>
            <el-icon><House /></el-icon>
            <span>房产管理</span>
          </template>
          <el-menu-item index="/realty">房产列表</el-menu-item>
          <el-menu-item v-if="hasRole(['GovernmentMSP'])" index="/realty/create">创建房产</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/transaction">
          <template #title>
            <el-icon><Sell /></el-icon>
            <span>交易管理</span>
          </template>
          <el-menu-item index="/transaction">交易列表</el-menu-item>
          <el-menu-item v-if="hasRole(['AgencyMSP', 'InvestorMSP'])" index="/transaction/create">创建交易</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/contract">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>合同管理</span>
          </template>
          <el-menu-item index="/contract">合同列表</el-menu-item>
          <el-menu-item v-if="hasRole(['AgencyMSP', 'InvestorMSP'])" index="/contract/create">创建合同</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/payment">
          <template #title>
            <el-icon><Money /></el-icon>
            <span>支付管理</span>
          </template>
          <el-menu-item index="/payment">支付列表</el-menu-item>
          <el-menu-item v-if="hasRole(['InvestorMSP', 'BankMSP'])" index="/payment/create">创建支付</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/tax">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>税费管理</span>
          </template>
          <el-menu-item index="/tax">税费列表</el-menu-item>
          <el-menu-item v-if="hasRole(['GovernmentMSP'])" index="/tax/create">创建税费</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/mortgage">
          <template #title>
            <el-icon><CreditCard /></el-icon>
            <span>抵押贷款</span>
          </template>
          <el-menu-item index="/mortgage">贷款列表</el-menu-item>
          <el-menu-item v-if="hasRole(['InvestorMSP'])" index="/mortgage/create">申请贷款</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu v-if="hasRole(['AdminMSP'])" index="/admin">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin/users">用户管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <!-- 主内容区 -->
    <el-container class="main-container">
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-icon class="toggle-sidebar" @click="toggleSidebar">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="$route.meta.title">{{ $route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <el-dropdown trigger="click" @command="handleCommand">
            <span class="user-dropdown">
              {{ username }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 内容区 -->
      <el-main class="content">
        <router-view />
      </el-main>
      
      <!-- 页脚 -->
      <el-footer class="footer">
        <div>房地产交易系统 &copy; {{ new Date().getFullYear() }}</div>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  Monitor, House, Sell, Document, Money, Wallet, 
  CreditCard, Setting, Fold, Expand, ArrowDown 
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 侧边栏折叠状态
const isCollapse = ref(false)

// 计算当前激活的菜单
const activeMenu = computed(() => {
  return router.currentRoute.value.path
})

// 用户名
const username = computed(() => userStore.username)

// 切换侧边栏折叠状态
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

// 检查角色权限
const hasRole = (roles) => {
  return userStore.hasRole(roles)
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push({ name: 'UserProfile' })
      break
    case 'changePassword':
      router.push({ name: 'ChangePassword' })
      break
    case 'logout':
      userStore.logout()
      break
  }
}
</script>

<style scoped>
.main-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  background-color: #001529;
  color: #fff;
  height: 100%;
  overflow-y: auto;
  transition: width 0.3s;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  color: #fff;
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
  background-color: #fff;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.toggle-sidebar {
  font-size: 20px;
  margin-right: 20px;
  cursor: pointer;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: #f5f7fa;
}

.footer {
  background-color: #fff;
  border-top: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  height: 40px;
  font-size: 12px;
  color: #999;
}
</style> 