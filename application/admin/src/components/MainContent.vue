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
  'investor': ['/realty', '/transaction', '/contract', '/payment', '/mortgage', '/chat'],
  'government': ['/realty', '/transaction', '/contract', '/tax', '/statistics', '/chat'],
  'bank': ['/transaction', '/payment', '/mortgage', '/statistics', '/contract', '/realty'],
  'audit': ['/contract', '/realty', '/transaction'],
  'thirdparty': ['/admin', '/realty', '/transaction', '/contract', '/payment', '/tax', '/mortgage', '/statistics']
}

// 检查当前路由是否有权限访问
const hasPermission = computed(() => {
  // 总是允许访问仪表盘和用户相关页面
  if (route.path === '/' || route.path === '/user/profile' || route.path === '/change-password') {
    return true
  }
  
  // 获取当前组织允许的路径
  const allowedPaths = organizationPermissions[userStore.organization] || []
  
  // 检查当前路径是否允许访问
  return allowedPaths.some(path => route.path.startsWith(path))
})

// 路由标题映射表
const routeTitleMap = {
  '': '首页',
  'realty': '房产管理',
  'transaction': '交易管理',
  'contract': '合同管理',
  'payment': '支付管理',
  'tax': '税费管理',
  'mortgage': '抵押贷款管理',
  'statistics': '统计分析',
  'user': '用户管理',
  'create': '创建',
  'edit': '编辑',
  'detail': '详情',
  'chat': '聊天室',
  'audit': '审核',
  'profile': '个人资料',
  'change-password': '修改密码',
  'transaction-statistics': '交易统计',
  'loan': '贷款统计'
}

// 面包屑导航
const breadcrumbs = computed(() => {
  const crumbs = [{ title: '首页', path: '/' }]
  
  if (route.path === '/') {
    return crumbs
  }
  
  const pathSegments = route.path.split('/').filter(Boolean)
  let currentPath = ''
  
  // 遍历路径段来构建面包屑
  for (let i = 0; i < pathSegments.length; i++) {
    const segment = pathSegments[i]
    currentPath += `/${segment}`
    
    // 查找当前路径匹配的路由配置
    const matchedRoute = router.resolve(currentPath).matched[0]
    
    // 查找该段路径对应的路由配置
    let title = ''
    
    // 如果是动态参数路径(如/:id)
    if (segment.match(/^[0-9a-fA-F]{8,}$/)) {
      // 看起来像UUID或ID，直接显示
      title = segment
    } else {
      // 优先从路由元数据获取标题
      if (matchedRoute && matchedRoute.meta && matchedRoute.meta.title) {
        // 当前路径有直接对应的路由配置
        title = matchedRoute.meta.title
      } else if (i > 0 && segment === 'create') {
        // 创建页面
        title = `创建${routeTitleMap[pathSegments[i-1].replace('-', '')] || ''}`
      } else if (i > 0 && segment === 'edit') {
        // 编辑页面
        title = `编辑${routeTitleMap[pathSegments[i-1].replace('-', '')] || ''}`
      } else if (i > 0 && segment === 'audit') {
        // 审核页面
        title = `${routeTitleMap[pathSegments[i-1].replace('-', '')] || ''}审核`
      } else {
        // 从映射表查找或直接使用段名
        title = routeTitleMap[segment] || segment
      }
    }
    
    // 添加到面包屑
    crumbs.push({
      title: title,
      path: currentPath
    })
  }
  
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