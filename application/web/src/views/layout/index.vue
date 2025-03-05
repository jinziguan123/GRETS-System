<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <el-menu
        :default-active="route.path"
        class="el-menu-vertical"
        router
      >
        <el-menu-item index="/">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/property">
          <el-icon><Office /></el-icon>
          <span>房产管理</span>
        </el-menu-item>
        <el-menu-item index="/transaction">
          <el-icon><Sell /></el-icon>
          <span>交易管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              {{ userStore.username }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { House, Office, Sell, ArrowDown } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const store = useStore()

const handleCommand = (command: string) => {
  if (command === 'logout') {
    store.dispatch('user/logout')
    router.push('/login')
  }
}
</script>

<style scoped>
.layout-container {
  height: 100%;
}

.el-menu-vertical {
  height: 100%;
}

.header-right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  height: 100%;
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
}
</style>