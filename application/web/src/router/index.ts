import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw, RouteMeta as VueRouteMeta } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 布局
const MainLayout = () => import('@/layouts/MainLayout.vue')
const AuthLayout = () => import('@/layouts/AuthLayout.vue')

// 认证页面
const Login = () => import('@/views/user/Login.vue')
const Register = () => import('@/views/user/Register.vue')

// 仪表盘
const Dashboard = () => import('@/views/dashboard/DashboardView.vue')

// 房产管理
const RealtyList = () => import('@/views/realty/RealtyList.vue')
const RealtyDetail = () => import('@/views/realty/RealtyDetail.vue')
const RealtyCreate = () => import('@/views/realty/RealtyCreate.vue')
const RealtyEdit = () => import('@/views/realty/RealtyEdit.vue')

// 交易管理
const TransactionList = () => import('@/views/transaction/TransactionList.vue')
const TransactionDetail = () => import('@/views/transaction/TransactionDetail.vue')
const TransactionCreate = () => import('@/views/transaction/TransactionCreate.vue')
// 新增交易聊天室
const TransactionChat = () => import('@/views/transaction/TransactionChat.vue')

// 合同管理
const ContractList = () => import('@/views/contract/ContractList.vue')
const ContractDetail = () => import('@/views/contract/ContractDetail.vue')
const ContractCreate = () => import('@/views/contract/ContractCreate.vue')
// 新增合同审核
const ContractAudit = () => import('@/views/contract/ContractAudit.vue')

// 支付管理
const PaymentList = () => import('@/views/payment/PaymentList.vue')
const PaymentDetail = () => import('@/views/payment/PaymentDetail.vue')
const PaymentCreate = () => import('@/views/payment/PaymentCreate.vue')

// 税费管理
const TaxList = () => import('@/views/tax/TaxList.vue')
const TaxDetail = () => import('@/views/tax/TaxDetail.vue')
const TaxCreate = () => import('@/views/tax/TaxCreate.vue')

// 抵押贷款管理
const MortgageList = () => import('@/views/mortgage/MortgageList.vue')
const MortgageDetail = () => import('@/views/mortgage/MortgageDetail.vue')
const MortgageCreate = () => import('@/views/mortgage/MortgageCreate.vue')
// 新增贷款审批
// const MortgageApprove = () => import('@/views/mortgage/MortgageApprove.vue')

// 统计分析
const TransactionStatistics = () => import('@/views/statistics/TransactionStatistics.vue')
const LoanStatistics = () => import('@/views/statistics/LoanStatistics.vue')

// 用户管理
const UserProfile = () => import('@/views/user/UserProfile.vue')
const ChangePassword = () => import('@/views/user/ChangePassword.vue')

// 管理员
// const AdminUserList = () => import('@/views/admin/UserList.vue')
// const AdminUserCreate = () => import('@/views/admin/UserCreate.vue')
// const AdminUserEdit = () => import('@/views/admin/UserEdit.vue')
// const AdminSystemSettings = () => import('@/views/admin/SystemSettings.vue')

interface RouteMeta extends VueRouteMeta {
  title?: string
  requiresAuth?: boolean
  organizations?: string[]
}

// 路由配置
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true } as RouteMeta,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard,
        meta: { title: '仪表盘' } as RouteMeta
      },
      // 房产管理路由
      {
        path: 'realty',
        name: 'RealtyList',
        component: RealtyList,
        meta: { title: '房产列表' } as RouteMeta
      },
      {
        path: 'realty/create',
        name: 'RealtyCreate',
        component: RealtyCreate,
        meta: { title: '创建房产', organizations: ['government'] } as RouteMeta
      },
      {
        path: 'realty/:id',
        name: 'RealtyDetail',
        component: RealtyDetail,
        meta: { title: '房产详情' } as RouteMeta
      },
      {
        path: 'realty/:id/edit',
        name: 'RealtyEdit',
        component: RealtyEdit,
        meta: { title: '编辑房产', organizations: ['government'] } as RouteMeta
      },
      // 交易管理路由
      {
        path: 'transaction',
        name: 'TransactionList',
        component: TransactionList,
        meta: { title: '交易列表', organizations: ['investor', 'bank', 'government'] } as RouteMeta
      },
      {
        path: 'transaction/create',
        name: 'TransactionCreate',
        component: TransactionCreate,
        meta: { title: '创建交易', organizations: ['investor'] } as RouteMeta
      },
      {
        path: 'transaction/:id',
        name: 'TransactionDetail',
        component: TransactionDetail,
        meta: { title: '交易详情', organizations: ['investor', 'bank', 'government', 'audit'] } as RouteMeta
      },
      {
        path: 'transaction/chat',
        name: 'TransactionChat',
        component: TransactionChat,
        meta: { title: '交易聊天室', organizations: ['investor'] } as RouteMeta
      },
      // 合同管理路由
      {
        path: 'contract',
        name: 'ContractList',
        component: ContractList,
        meta: { title: '合同列表', organizations: ['investor', 'bank', 'government', 'audit'] } as RouteMeta
      },
      {
        path: 'contract/create',
        name: 'ContractCreate',
        component: ContractCreate,
        meta: { title: '创建合同', organizations: ['investor', 'government'] } as RouteMeta
      },
      {
        path: 'contract/:id',
        name: 'ContractDetail',
        component: ContractDetail,
        meta: { title: '合同详情', organizations: ['investor', 'bank', 'government', 'audit'] } as RouteMeta
      },
      {
        path: 'contract/audit',
        name: 'ContractAudit',
        component: ContractAudit,
        meta: { title: '合同审核', organizations: ['audit'] } as RouteMeta
      },
      // 支付管理路由
      {
        path: 'payment',
        name: 'PaymentList',
        component: PaymentList,
        meta: { title: '支付列表', organizations: ['investor', 'bank'] } as RouteMeta
      },
      {
        path: 'payment/create',
        name: 'PaymentCreate',
        component: PaymentCreate,
        meta: { title: '创建支付', organizations: ['investor'] } as RouteMeta
      },
      {
        path: 'payment/:id',
        name: 'PaymentDetail',
        component: PaymentDetail,
        meta: { title: '支付详情', organizations: ['investor', 'bank'] } as RouteMeta
      },
      // 税费管理路由
      {
        path: 'tax',
        name: 'TaxList',
        component: TaxList,
        meta: { title: '税费列表', organizations: ['government', 'investor'] } as RouteMeta
      },
      {
        path: 'tax/create',
        name: 'TaxCreate',
        component: TaxCreate,
        meta: { title: '创建税费', organizations: ['government'] } as RouteMeta
      },
      {
        path: 'tax/:id',
        name: 'TaxDetail',
        component: TaxDetail,
        meta: { title: '税费详情', organizations: ['government', 'investor'] } as RouteMeta
      },
      // 抵押贷款管理路由
      {
        path: 'mortgage',
        name: 'MortgageList',
        component: MortgageList,
        meta: { title: '抵押贷款列表', organizations: ['investor', 'bank'] } as RouteMeta
      },
      {
        path: 'mortgage/create',
        name: 'MortgageCreate',
        component: MortgageCreate,
        meta: { title: '创建抵押贷款', organizations: ['investor'] } as RouteMeta
      },
      {
        path: 'mortgage/:id',
        name: 'MortgageDetail',
        component: MortgageDetail,
        meta: { title: '抵押贷款详情', organizations: ['investor', 'bank'] } as RouteMeta
      },
      // 统计分析路由
      {
        path: 'statistics/transaction',
        name: 'TransactionStatistics',
        component: TransactionStatistics,
        meta: { title: '交易统计', organizations: ['government'] } as RouteMeta
      },
      {
        path: 'statistics/loan',
        name: 'LoanStatistics',
        component: LoanStatistics,
        meta: { title: '贷款统计', organizations: ['bank'] } as RouteMeta
      },
      // 用户管理路由
      {
        path: 'profile',
        name: 'UserProfile',
        component: UserProfile,
        meta: { title: '个人资料' } as RouteMeta
      },
      {
        path: 'change-password',
        name: 'ChangePassword',
        component: ChangePassword,
        meta: { title: '修改密码' } as RouteMeta
      }
    ]
  },
  {
    path: '/auth',
    component: AuthLayout,
    meta: { requiresAuth: false } as RouteMeta,
    children: [
      {
        path: 'login',
        name: 'Login',
        component: Login,
        meta: { title: '登录' } as RouteMeta
      },
      {
        path: 'register',
        name: 'Register',
        component: Register,
        meta: { title: '注册' } as RouteMeta
      }
    ]
  },
  // 404页面
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, _, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 房地产交易系统` : '房地产交易系统'
  
  // 获取token和用户存储
  const userStore = useUserStore()
  const token = localStorage.getItem('token')
  
  // 如果有token但没有设置axios头部，初始化头部
  if (token && !axios.defaults.headers.common['Authorization']) {
    axios.defaults.headers.common['Authorization'] = `${token}`
  }
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !token) {
    // 需要认证但没有token，重定向到登录页
    ElMessage.warning('请先登录后再访问该页面')
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }
  
  // 如果去登录页但已经登录，重定向到首页
  if ((to.name === 'Login' || to.name === 'Register') && token) {
    next({ name: 'Dashboard' })
    return
  }
  
  // 检查组织权限
  const organizations = to.meta.organizations as string[] | undefined
  if (organizations && organizations.length > 0 && token) {
    // 获取用户组织
    const userOrganization = userStore.organization
    
    // 如果用户组织不在允许的组织列表中，重定向到仪表盘
    if (!organizations.includes(userOrganization)) {
      ElMessage.warning(`当前页面仅对${organizations.map((org: string) => userStore.getOrganizationName(org)).join('、')}开放访问`)
      next({ name: 'Dashboard' })
      return
    }
  }
  
  // 允许访问
  next()
})

export default router 