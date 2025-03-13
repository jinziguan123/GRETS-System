<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <div class="logo">
            <h1>GRETS系统</h1>
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

<script setup>
// 认证布局组件，用于登录、注册等页面
</script>

<style scoped>
.auth-layout {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1976d2, #2196f3);
  padding: 20px;
}

.auth-container {
  width: 100%;
  max-width: 440px;
}

.auth-card {
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
</style> 