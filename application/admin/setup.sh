#!/bin/bash

echo "======== 房地产交易系统前端工具 ========"
echo "当前Node版本:"
node -v
echo "当前NPM版本:"
npm -v

# 安装依赖
install_deps() {
  echo "=== 安装依赖 ==="
  
  # 安装Vue Router
  npm install vue-router@4
  
  # 安装Pinia状态管理
  npm install pinia
  
  # 安装Element Plus
  npm install element-plus
  
  # 安装Element Plus图标
  npm install @element-plus/icons-vue
  
  # 安装Axios
  npm install axios
  
  # 安装VueUse
  npm install @vueuse/core
  
  # 安装SASS
  npm install -D sass
  
  # 安装ECharts
  npm install echarts
  
  # 安装dayjs处理日期
  npm install dayjs
  
  # 安装unplugin-auto-import和unplugin-vue-components
  npm install -D unplugin-auto-import unplugin-vue-components
  
  echo "依赖安装完成"
}

# 创建目录结构
create_structure() {
  echo "=== 创建目录结构 ==="
  
  # 创建路由目录
  mkdir -p src/router
  
  # 创建状态管理目录
  mkdir -p src/stores
  
  # 创建视图目录
  mkdir -p src/views/dashboard
  mkdir -p src/views/realty
  mkdir -p src/views/transaction
  mkdir -p src/views/contract
  mkdir -p src/views/payment
  mkdir -p src/views/tax
  mkdir -p src/views/mortgage
  mkdir -p src/views/user
  mkdir -p src/views/admin
  
  # 创建API目录
  mkdir -p src/api
  
  # 创建工具目录
  mkdir -p src/utils
  
  # 创建布局目录
  mkdir -p src/layouts
  
  # 创建组件目录
  mkdir -p src/components/common
  mkdir -p src/components/realty
  mkdir -p src/components/transaction
  
  echo "目录结构创建完成"
}

# 运行开发服务器
run_dev() {
  echo "=== 运行开发服务器 ==="
  npm run dev
}

# 构建生产版本
build() {
  echo "=== 构建生产版本 ==="
  npm run build
}

# 帮助信息
show_help() {
  echo "用法: ./setup.sh [命令]"
  echo "命令:"
  echo "  deps     - 安装依赖"
  echo "  struct   - 创建目录结构"
  echo "  dev      - 运行开发服务器"
  echo "  build    - 构建生产版本"
  echo "  help     - 显示帮助信息"
}

# 主函数
main() {
  if [ $# -eq 0 ]; then
    show_help
    exit 0
  fi

  case "$1" in
    "deps")
      install_deps
      ;;
    "struct")
      create_structure
      ;;
    "dev")
      run_dev
      ;;
    "build")
      build
      ;;
    "help")
      show_help
      ;;
    *)
      echo "未知命令: $1"
      show_help
      exit 1
      ;;
  esac
}

main "$@" 