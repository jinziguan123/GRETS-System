<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-logo">
        <h1>GRETS</h1>
        <p>房地产交易系统</p>
      </div>
      
      <div class="auth-content">
        <router-view />
      </div>
      
      <div class="auth-footer">
        <p>房地产交易系统 &copy; {{ new Date().getFullYear() }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
// 认证布局不需要特殊逻辑
</script>

<style scoped>
.auth-layout {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.auth-container {
  width: 400px;
  padding: 40px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.auth-logo {
  text-align: center;
  margin-bottom: 30px;
}

.auth-logo h1 {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
  margin: 0;
}

.auth-logo p {
  font-size: 16px;
  color: #606266;
  margin: 10px 0 0;
}

.auth-content {
  margin-bottom: 30px;
}

.auth-footer {
  text-align: center;
  font-size: 12px;
  color: #909399;
}
</style> 