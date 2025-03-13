// 模拟API响应数据，用于开发测试

// 仪表盘数据
export const dashboardMock = {
  code: 200,
  data: {
    statistics: {
      realtyCount: 127,
      transactionCount: 85,
      contractCount: 94,
      taxAmount: 1358000
    },
    latestTransactions: [
      { id: 'T-1001', realtyName: '翠湖豪庭 3号楼605', amount: 1580000, status: '完成' },
      { id: 'T-1002', realtyName: '阳光花园 12栋1203', amount: 2450000, status: '进行中' },
      { id: 'T-1003', realtyName: '蓝天公寓 B座502', amount: 1260000, status: '完成' },
      { id: 'T-1004', realtyName: '江南名府 6号楼1801', amount: 3280000, status: '待处理' }
    ],
    latestContracts: [
      { id: 'C-2001', title: '翠湖豪庭605购房合同', date: '2024-05-12', status: '完成' },
      { id: 'C-2002', title: '阳光花园1203认购协议', date: '2024-05-10', status: '进行中' },
      { id: 'C-2003', title: '蓝天公寓502过户协议', date: '2024-05-08', status: '完成' },
      { id: 'C-2004', title: '江南名府1801租赁合同', date: '2024-05-05', status: '待处理' }
    ]
  },
  message: '获取仪表盘数据成功'
}

// 用户数据
export const userMock = {
  login: {
    code: 200,
    data: {
      token: 'mock-token-123456789',
      username: 'testuser',
      userRole: 'InvestorMSP',
      userId: 'user-001'
    },
    message: '登录成功'
  },
  userInfo: {
    code: 200,
    data: {
      id: 'user-001',
      username: 'testuser',
      fullName: '张三',
      email: 'test@example.com',
      phone: '13800138000',
      role: 'InvestorMSP',
      status: 'active',
      avatar: '',
      createTime: '2024-01-01'
    },
    message: '获取用户信息成功'
  }
}

// 房产数据
export const realtyMock = {
  list: {
    code: 200,
    data: [
      {
        id: 'R001',
        title: '阳光花园 3室2厅',
        type: '住宅',
        address: '杭州市西湖区文三路138号阳光花园3幢1单元601',
        area: 120,
        price: 1200000,
        status: 'available',
        owner: 'user-001',
        createTime: '2024-01-10'
      },
      {
        id: 'R002',
        title: '江南名府 4室2厅',
        type: '住宅',
        address: '杭州市滨江区滨盛路1509号江南名府12幢2单元1801',
        area: 180,
        price: 2500000,
        status: 'available',
        owner: 'user-001',
        createTime: '2024-02-15'
      },
      {
        id: 'R003',
        title: '城市广场 商铺A12',
        type: '商铺',
        address: '杭州市江干区凯旋路166号城市广场A区12号',
        area: 85,
        price: 3500000,
        status: 'available',
        owner: 'user-001',
        createTime: '2024-03-20'
      }
    ],
    total: 3,
    message: '获取房产列表成功'
  },
  detail: {
    code: 200,
    data: {
      id: 'R001',
      title: '阳光花园 3室2厅',
      type: '住宅',
      address: '杭州市西湖区文三路138号阳光花园3幢1单元601',
      area: 120,
      price: 1200000,
      buildYear: 2015,
      floorInfo: '6/18',
      orientation: '南',
      decoration: '精装修',
      status: 'available',
      facilities: ['电梯', '车位', '天然气', '暖气', '宽带'],
      description: '房子位于市中心地段，交通便利，周边设施齐全，采光好，户型方正。',
      owner: {
        id: 'user-001',
        name: '张三',
        phone: '13800138000'
      },
      images: [
        'https://example.com/image1.jpg',
        'https://example.com/image2.jpg',
        'https://example.com/image3.jpg'
      ],
      documents: [
        {
          id: 'doc-001',
          name: '房产证',
          url: 'https://example.com/doc1.pdf'
        }
      ],
      createTime: '2024-01-10',
      updateTime: '2024-01-10'
    },
    message: '获取房产详情成功'
  }
}

// 交易数据
export const transactionMock = {
  list: {
    code: 200,
    data: [
      {
        id: 'T-1001',
        realty: {
          id: 'R001',
          title: '阳光花园 3室2厅'
        },
        buyer: {
          id: 'user-002',
          name: '李四'
        },
        seller: {
          id: 'user-001',
          name: '张三'
        },
        amount: 1580000,
        status: '完成',
        createTime: '2024-04-01',
        completeTime: '2024-04-15'
      },
      {
        id: 'T-1002',
        realty: {
          id: 'R002',
          title: '江南名府 4室2厅'
        },
        buyer: {
          id: 'user-003',
          name: '王五'
        },
        seller: {
          id: 'user-001',
          name: '张三'
        },
        amount: 2450000,
        status: '进行中',
        createTime: '2024-05-01'
      }
    ],
    total: 2,
    message: '获取交易列表成功'
  },
  detail: {
    code: 200,
    data: {
      id: 'T-1001',
      realty: {
        id: 'R001',
        title: '阳光花园 3室2厅',
        address: '杭州市西湖区文三路138号阳光花园3幢1单元601',
        area: 120,
        price: 1200000
      },
      buyer: {
        id: 'user-002',
        name: '李四',
        phone: '13900139000',
        email: 'lisi@example.com'
      },
      seller: {
        id: 'user-001',
        name: '张三',
        phone: '13800138000',
        email: 'zhangsan@example.com'
      },
      agent: {
        id: 'user-004',
        name: '赵六',
        phone: '13700137000',
        email: 'zhaoliu@example.com',
        agency: '好房中介'
      },
      amount: 1580000,
      downPayment: 300000,
      status: '完成',
      steps: [
        { name: '创建交易', status: '完成', time: '2024-04-01' },
        { name: '签订合同', status: '完成', time: '2024-04-05' },
        { name: '支付定金', status: '完成', time: '2024-04-08' },
        { name: '办理贷款', status: '完成', time: '2024-04-10' },
        { name: '支付尾款', status: '完成', time: '2024-04-12' },
        { name: '过户登记', status: '完成', time: '2024-04-15' }
      ],
      documents: [
        { id: 'doc-101', name: '购房合同', type: 'contract' },
        { id: 'doc-102', name: '房产证明', type: 'certificate' },
        { id: 'doc-103', name: '贷款合同', type: 'mortgage' }
      ],
      comments: [
        { time: '2024-04-02', user: '赵六', content: '已联系买卖双方，准备签订合同' },
        { time: '2024-04-06', user: '赵六', content: '合同已签订，准备办理贷款手续' },
        { time: '2024-04-16', user: '赵六', content: '交易已完成，房产证变更手续已办理' }
      ],
      createTime: '2024-04-01',
      updateTime: '2024-04-15',
      completeTime: '2024-04-15'
    },
    message: '获取交易详情成功'
  }
}

// 合同数据
export const contractMock = {
  list: {
    code: 200,
    data: [
      {
        id: 'C-2001',
        title: '翠湖豪庭605购房合同',
        type: '购房合同',
        transaction: {
          id: 'T-1001',
          realty: '阳光花园 3室2厅'
        },
        parties: [
          { name: '张三', role: '卖方' },
          { name: '李四', role: '买方' }
        ],
        status: '已签署',
        date: '2024-04-05',
        expireDate: '2024-07-05'
      },
      {
        id: 'C-2002',
        title: '阳光花园1203认购协议',
        type: '认购协议',
        transaction: {
          id: 'T-1002',
          realty: '江南名府 4室2厅'
        },
        parties: [
          { name: '张三', role: '卖方' },
          { name: '王五', role: '买方' }
        ],
        status: '进行中',
        date: '2024-05-10',
        expireDate: '2024-08-10'
      }
    ],
    total: 2,
    message: '获取合同列表成功'
  }
}

// 抵押贷款数据
export const mortgageMock = {
  list: {
    code: 200,
    data: [
      {
        id: 'M-001',
        realty: {
          id: 'R001',
          title: '阳光花园 3室2厅'
        },
        applicant: '李四',
        amount: 1000000,
        term: 20,
        interestRate: 4.9,
        bank: '中国建设银行',
        status: '已审批',
        applyDate: '2024-04-08',
        approveDate: '2024-04-10'
      }
    ],
    total: 1,
    message: '获取贷款列表成功'
  },
  detail: {
    code: 200,
    data: {
      id: 'M-001',
      realty: {
        id: 'R001',
        title: '阳光花园 3室2厅',
        address: '杭州市西湖区文三路138号阳光花园3幢1单元601',
        price: 1200000
      },
      applicant: {
        name: '李四',
        idCard: '330102198801010011',
        phone: '13900139000',
        email: 'lisi@example.com',
        employer: 'ABC科技公司',
        income: 300000
      },
      amount: 1000000,
      term: 20,
      rateType: 'fixed',
      interestRate: 4.9,
      bank: {
        id: 'bank2',
        name: '中国建设银行'
      },
      repaymentMethod: 'equal_installment',
      repaymentAccount: '6227001234567890123',
      monthlyPayment: 6552.18,
      status: '已审批',
      documents: [
        { id: 'file-001', name: '身份证复印件' },
        { id: 'file-002', name: '收入证明' }
      ],
      reviewComments: '申请人信用良好，房产价值和变现能力高，批准贷款申请。',
      applyDate: '2024-04-08',
      approveDate: '2024-04-10',
      firstPaymentDate: '2024-05-10'
    },
    message: '获取贷款详情成功'
  }
}

// 导出所有模拟数据
export default {
  dashboardMock,
  userMock,
  realtyMock,
  transactionMock,
  contractMock,
  mortgageMock
} 