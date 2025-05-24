<template>
  <div>
    <router-view />
  </div>
</template>

<script setup lang="ts">
// App.vue 是应用的根组件，只需要包含路由视图
</script>

<style>
/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100vh;
  width: 100vw;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-size: 14px;
  color: #333;
  background-color: #f5f7fa;
}

#app {
  height: 100vh;
  width: 100vw;
}

/* 常用辅助类 */
.text-center {
  text-align: center;
}

.mt-10 {
  margin-top: 10px;
}

.mb-10 {
  margin-bottom: 10px;
}

.ml-10 {
  margin-left: 10px;
}

.mr-10 {
  margin-right: 10px;
}

.p-20 {
  padding: 20px;
}

/* 卡片样式 */
.card {
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.card-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

/* 表格工具栏 */
.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.search-form {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

/* 详情页样式 */
.detail-item {
  margin-bottom: 15px;
}

.detail-label {
  font-weight: bold;
  color: #606266;
  margin-bottom: 5px;
}

.detail-value {
  color: #303133;
}

/* 表单样式 */
.form-container {
  max-width: 800px;
  margin: 0 auto;
}

.form-footer {
  margin-top: 30px;
  text-align: center;
}

.form-footer .el-button {
  min-width: 120px;
}
</style>
