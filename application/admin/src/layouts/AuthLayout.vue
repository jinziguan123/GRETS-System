<template>
  <div class="auth-layout">
    <!-- 左侧区域 - 品牌信息 -->
    <div class="brand-section">
      <div class="brand-content">
        <h1 class="brand-title">GRETS<br/>系统后台管理</h1>
        <p class="brand-description">基于区块链的房<br/>地产交易系统</p>
        <!-- 可以在这里添加品牌图片或其他装饰元素 -->
      </div>
    </div>

    <!-- 右侧区域 - 登录表单 -->
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <div class="logo">
            <h1>GRETS系统后台管理</h1>
            <p>基于区块链的房地产交易系统</p>
          </div>
        </div>

        <div class="auth-content">
          <router-view />
        </div>

        <div class="auth-footer">
          <p>&copy; {{ new Date().getFullYear() }} GRETS System. 版权所有</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 认证布局组件，用于登录、注册等页面
</script>

<style scoped>
.auth-layout {
  display: flex;
  width: 100vw;  /* 占据整个视口宽度 */
  height: 100vh; /* 占据整个视口高度 */
  overflow: hidden; /* 防止可能的滚动条 */
  background: linear-gradient(135deg, #1976d2, #2196f3);
}

/* 左侧品牌区域 */
.brand-section {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
  padding: 40px;
  text-align: center;
}

.brand-content {
  max-width: 600px;
}

.brand-title {
  font-size: 56px;
  font-weight: bold;
  margin-bottom: 20px;
  line-height: 1.2;
}

.brand-description {
  font-size: 24px;
  line-height: 1.5;
}

/* 右侧认证区域 */
.auth-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.auth-card {
  width: 100%;
  max-width: 440px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.auth-header {
  padding: 30px 30px 20px;
  text-align: center;
}

.logo h1 {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 8px;
  color: #1976d2;
}

.logo p {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.auth-content {
  padding: 0 30px 30px;
}

.auth-footer {
  padding: 15px 0;
  text-align: center;
  background-color: #f5f7fa;
  border-top: 1px solid #eee;
}

.auth-footer p {
  margin: 0;
  font-size: 12px;
  color: #999;
}

/* 移动设备上的响应式布局 */
@media (max-width: 992px) {
  .auth-layout {
    flex-direction: column;
  }

  .brand-section {
    padding: 40px 20px;
    min-height: 40vh;
  }

  .auth-container {
    min-height: 60vh;
  }

  .auth-card {
    max-width: 100%;
  }
}

/* 额外确保页面无边距 */
:deep(body), :deep(html) {
  margin: 0;
  padding: 0;
  overflow: hidden;
}
</style> 