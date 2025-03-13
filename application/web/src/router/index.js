import { createRouter, createWebHistory } from 'vue-router'

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

// 合同管理
const ContractList = () => import('@/views/contract/ContractList.vue')
const ContractDetail = () => import('@/views/contract/ContractDetail.vue')
const ContractCreate = () => import('@/views/contract/ContractCreate.vue')

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

// 用户管理
const UserProfile = () => import('@/views/user/UserProfile.vue')
const ChangePassword = () => import('@/views/user/ChangePassword.vue')

// 管理员
const AdminUserList = () => import('@/views/admin/UserList.vue')
const AdminUserCreate = () => import('@/views/admin/UserCreate.vue')
const AdminUserEdit = () => import('@/views/admin/UserEdit.vue')

// 路由配置
const routes = [
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard,
        meta: { title: '仪表盘' }
      },
      // 房产管理路由
      {
        path: 'realty',
        name: 'RealtyList',
        component: RealtyList,
        meta: { title: '房产列表' }
      },
      {
        path: 'realty/create',
        name: 'RealtyCreate',
        component: RealtyCreate,
        meta: { title: '创建房产', roles: ['GovernmentMSP'] }
      },
      {
        path: 'realty/:id',
        name: 'RealtyDetail',
        component: RealtyDetail,
        meta: { title: '房产详情' }
      },
      {
        path: 'realty/:id/edit',
        name: 'RealtyEdit',
        component: RealtyEdit,
        meta: { title: '编辑房产', roles: ['GovernmentMSP'] }
      },
      // 交易管理路由
      {
        path: 'transaction',
        name: 'TransactionList',
        component: TransactionList,
        meta: { title: '交易列表' }
      },
      {
        path: 'transaction/create',
        name: 'TransactionCreate',
        component: TransactionCreate,
        meta: { title: '创建交易', roles: ['AgencyMSP', 'InvestorMSP'] }
      },
      {
        path: 'transaction/:id',
        name: 'TransactionDetail',
        component: TransactionDetail,
        meta: { title: '交易详情' }
      },
      // 合同管理路由
      {
        path: 'contract',
        name: 'ContractList',
        component: ContractList,
        meta: { title: '合同列表' }
      },
      {
        path: 'contract/create',
        name: 'ContractCreate',
        component: ContractCreate,
        meta: { title: '创建合同', roles: ['AgencyMSP', 'InvestorMSP'] }
      },
      {
        path: 'contract/:id',
        name: 'ContractDetail',
        component: ContractDetail,
        meta: { title: '合同详情' }
      },
      // 支付管理路由
      {
        path: 'payment',
        name: 'PaymentList',
        component: PaymentList,
        meta: { title: '支付列表' }
      },
      {
        path: 'payment/create',
        name: 'PaymentCreate',
        component: PaymentCreate,
        meta: { title: '创建支付', roles: ['InvestorMSP', 'BankMSP'] }
      },
      {
        path: 'payment/:id',
        name: 'PaymentDetail',
        component: PaymentDetail,
        meta: { title: '支付详情' }
      },
      // 税费管理路由
      {
        path: 'tax',
        name: 'TaxList',
        component: TaxList,
        meta: { title: '税费列表' }
      },
      {
        path: 'tax/create',
        name: 'TaxCreate',
        component: TaxCreate,
        meta: { title: '创建税费', roles: ['GovernmentMSP'] }
      },
      {
        path: 'tax/:id',
        name: 'TaxDetail',
        component: TaxDetail,
        meta: { title: '税费详情' }
      },
      // 抵押贷款管理路由
      {
        path: 'mortgage',
        name: 'MortgageList',
        component: MortgageList,
        meta: { title: '抵押贷款列表' }
      },
      {
        path: 'mortgage/create',
        name: 'MortgageCreate',
        component: MortgageCreate,
        meta: { title: '创建抵押贷款', roles: ['InvestorMSP'] }
      },
      {
        path: 'mortgage/:id',
        name: 'MortgageDetail',
        component: MortgageDetail,
        meta: { title: '抵押贷款详情' }
      },
      // 用户管理路由
      {
        path: 'profile',
        name: 'UserProfile',
        component: UserProfile,
        meta: { title: '个人资料' }
      },
      {
        path: 'change-password',
        name: 'ChangePassword',
        component: ChangePassword,
        meta: { title: '修改密码' }
      },
      // 管理员路由
      {
        path: 'admin/users',
        name: 'AdminUserList',
        component: AdminUserList,
        meta: { title: '用户管理', roles: ['AdminMSP'] }
      },
      {
        path: 'admin/users/create',
        name: 'AdminUserCreate',
        component: AdminUserCreate,
        meta: { title: '创建用户', roles: ['AdminMSP'] }
      },
      {
        path: 'admin/users/:id/edit',
        name: 'AdminUserEdit',
        component: AdminUserEdit,
        meta: { title: '编辑用户', roles: ['AdminMSP'] }
      }
    ]
  },
  {
    path: '/auth',
    component: AuthLayout,
    meta: { requiresAuth: false },
    children: [
      {
        path: 'login',
        name: 'Login',
        component: Login,
        meta: { title: '登录' }
      },
      {
        path: 'register',
        name: 'Register',
        component: Register,
        meta: { title: '注册' }
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
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 房地产交易系统` : '房地产交易系统'
  
  // 获取token
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('userRole')
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !token) {
    // 需要认证但没有token，重定向到登录页
    next({ name: 'Login' })
  } 
  // 检查角色权限
  else if (to.meta.roles && !to.meta.roles.includes(userRole)) {
    // 没有权限访问该页面
    next({ name: 'Dashboard' })
  } 
  else {
    // 允许访问
    next()
  }
})

export default router 