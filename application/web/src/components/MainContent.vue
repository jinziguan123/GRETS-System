<template>
  <div class="main-content-wrapper">
    <!-- 面包屑导航 -->
    <el-breadcrumb class="breadcrumb" separator="/">
      <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path">
        {{ item.title }}
      </el-breadcrumb-item>
    </el-breadcrumb>
    
    <!-- 权限不足提示 -->
    <el-result
      v-if="!hasPermission"
      icon="warning"
      title="权限不足"
      sub-title="您没有权限访问当前页面"
    >
      <template #extra>
        <el-button type="primary" @click="goToDashboard">返回首页</el-button>
      </template>
    </el-result>
    
    <!-- 内容区域 -->
    <div v-else class="content-container">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 组织权限映射
const organizationPermissions = {
  'investor': ['/realty', '/transaction', '/contract', '/payment', '/mortgage'],
  'government': ['/realty', '/transaction', '/contract', '/tax', '/statistics'],
  'bank': ['/transaction', '/payment', '/mortgage', '/statistics'],
  'audit': ['/contract', '/realty'],
  'administrator': ['/admin', '/realty', '/transaction', '/contract', '/payment', '/tax', '/mortgage', '/statistics']
}

// 检查当前路由是否有权限访问
const hasPermission = computed(() => {
  // 总是允许访问仪表盘和用户相关页面
  if (route.path === '/' || route.path === '/profile' || route.path === '/change-password') {
    return true
  }
  
  // 获取当前组织允许的路径
  const allowedPaths = organizationPermissions[userStore.organization] || []
  
  // 检查当前路径是否允许访问
  return allowedPaths.some(path => route.path.startsWith(path))
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

// 返回首页
const goToDashboard = () => {
  router.push('/')
}
</script>

<style scoped>
.main-content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.breadcrumb {
  margin-bottom: 20px;
}

.content-container {
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  min-height: calc(100% - 40px);
  flex: 1;
  padding: 20px;
}

@media (max-width: 768px) {
  .content-container {
    padding: 15px;
  }
}
</style> 