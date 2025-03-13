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
          <el-menu-item index="/realty/list">房产列表</el-menu-item>
          <el-menu-item index="/realty/create" v-if="hasPermission(['GovernmentMSP'])">添加房产</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/transaction">
          <template #title>
            <el-icon><Sell /></el-icon>
            <span>交易管理</span>
          </template>
          <el-menu-item index="/transaction/list">交易列表</el-menu-item>
          <el-menu-item index="/transaction/create" v-if="hasPermission(['AgencyMSP', 'InvestorMSP'])">创建交易</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/contract">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>合同管理</span>
          </template>
          <el-menu-item index="/contract/list">合同列表</el-menu-item>
          <el-menu-item index="/contract/create" v-if="hasPermission(['AgencyMSP', 'InvestorMSP'])">创建合同</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/payment">
          <template #title>
            <el-icon><Money /></el-icon>
            <span>支付管理</span>
          </template>
          <el-menu-item index="/payment/list">支付列表</el-menu-item>
          <el-menu-item index="/payment/create" v-if="hasPermission(['InvestorMSP', 'BankMSP'])">创建支付</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/tax">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>税费管理</span>
          </template>
          <el-menu-item index="/tax/list">税费列表</el-menu-item>
          <el-menu-item index="/tax/create" v-if="hasPermission(['GovernmentMSP'])">创建税费</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/mortgage">
          <template #title>
            <el-icon><CreditCard /></el-icon>
            <span>抵押贷款</span>
          </template>
          <el-menu-item index="/mortgage/list">贷款列表</el-menu-item>
          <el-menu-item index="/mortgage/create" v-if="hasPermission(['InvestorMSP'])">申请贷款</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/admin" v-if="hasPermission(['AdminMSP'])">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/admin/user">用户管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <!-- 主要内容区域 -->
    <el-container class="main-container">
      <el-header class="header">
        <div class="header-left">
          <el-button type="text" @click="toggleSidebar">
            <el-icon><Menu /></el-icon>
          </el-button>
        </div>
        <div class="header-right">
          <el-dropdown trigger="click">
            <span class="user-dropdown">
              {{ username }}
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/user/profile')">个人信息</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/user/password')">修改密码</el-dropdown-item>
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <!-- 面包屑导航 -->
        <el-breadcrumb class="breadcrumb" separator="/">
          <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path">
            {{ item.title }}
          </el-breadcrumb-item>
        </el-breadcrumb>
        
        <!-- 内容区域 -->
        <div class="content-container">
          <router-view />
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  Menu, 
  House, 
  Document, 
  Sell, 
  Money, 
  Wallet, 
  CreditCard, 
  Setting, 
  Monitor,
  ArrowDown
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 用户信息
const username = computed(() => userStore.username)
const userRole = computed(() => userStore.userRole)

// 侧边栏折叠状态
const isCollapse = ref(false)
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

// 当前激活的菜单项
const activeMenu = computed(() => {
  return route.path
})

// 面包屑导航
const breadcrumbs = computed(() => {
  const crumbs = [{ title: '首页', path: '/' }]
  
  if (route.path === '/') {
    return crumbs
  }
  
  const paths = route.path.split('/').filter(Boolean)
  
  let currentPath = ''
  paths.forEach(path => {
    currentPath += `/${path}`
    const currentRoute = router.resolve(currentPath).matched[0]
    if (currentRoute) {
      crumbs.push({
        title: currentRoute.meta.title || path,
        path: currentPath
      })
    }
  })
  
  return crumbs
})

// 权限检查
const hasPermission = (roles) => {
  if (!roles || roles.length === 0) return true
  return roles.includes(userRole.value)
}

// 退出登录
const logout = () => {
  userStore.logout()
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

.user-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
  font-size: 14px;
}

.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 20px;
}

.content-container {
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  min-height: calc(100% - 40px);
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
  
  .content-container {
    padding: 15px;
  }
}
</style> 