<template>
  <div class="realty-list">
    <div class="page-header">
      <h2>房产列表</h2>
      <el-button 
        type="primary" 
        v-if="userStore.hasOrganization('government')"
        @click="router.push('/realty/create')"
      >
        添加房产
      </el-button>
    </div>

    <!-- 搜索条件 -->
    <el-card class="filter-container">
      <el-form :model="searchForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="证书号">
              <el-input v-model="searchForm.realtyCert" placeholder="请输入不动产证号" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房产类型">
              <el-select v-model="searchForm.realtyType" placeholder="请选择房产类型" clearable style="width: 100%">
                <el-option label="住宅" value="HOUSE"></el-option>
                <el-option label="商铺" value="SHOP"></el-option>
                <el-option label="办公" value="OFFICE"></el-option>
                <el-option label="工业" value="INDUSTRIAL"></el-option>
                <el-option label="其他" value="OTHER"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="户型">
              <el-select v-model="searchForm.houseType" placeholder="请选择户型" clearable style="width: 100%">
                <el-option label="一室" value="single"></el-option>
                <el-option label="两室" value="double"></el-option>
                <el-option label="三室" value="triple"></el-option>
                <el-option label="四室及以上" value="multiple"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="省份">
              <el-select 
                v-model="searchForm.province" 
                placeholder="请选择省份" 
                clearable 
                style="width: 100%"
                @change="handleProvinceChange"
                filterable
              >
                <el-option 
                  v-for="province in provinceList" 
                  :key="province.provinceCode" 
                  :label="province.provinceName" 
                  :value="province.provinceName"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="城市">
              <el-select 
                v-model="searchForm.city" 
                placeholder="请选择城市" 
                clearable 
                style="width: 100%"
                :disabled="!searchForm.province"
                @change="handleCityChange"
                filterable
              >
                <el-option 
                  v-for="city in cityList" 
                  :key="city" 
                  :label="city" 
                  :value="city"
                ></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="区/县">
              <el-input v-model="searchForm.district" placeholder="请输入区/县" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="街道">
              <el-input v-model="searchForm.street" placeholder="请输入街道" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="小区">
              <el-input v-model="searchForm.community" placeholder="请输入小区" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="单元">
              <el-input v-model="searchForm.unit" placeholder="请输入单元" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="楼层">
              <el-input v-model="searchForm.floor" placeholder="请输入楼层" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房号">
              <el-input v-model="searchForm.room" placeholder="请输入房号" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="是否为新房">
              <el-select v-model="searchForm.isNewHouse" placeholder="请选择是否为新房" clearable>
                <el-option label="是" :value="true"></el-option>
                <el-option label="否" :value="false"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="房产状态">
              <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 100%">
                <el-option label="正常" value="NORMAL"></el-option>
                <el-option label="挂牌" value="PENDING_SALE"></el-option>
                <el-option label="已抵押" value="IN_MORTGAGE"></el-option>
                <el-option label="已冻结" value="FROZEN"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="面积范围">
              <el-input-number v-model="searchForm.minArea" :min="0.0" :step="10" placeholder="最小面积" style="width: 45%" />
              <span class="separator">-</span>
              <el-input-number v-model="searchForm.maxArea" :min="0.0" :step="10" placeholder="最大面积" style="width: 45%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="价格范围">
              <el-input-number v-model="searchForm.minPrice" :min="0.0" :step="50000" :formatter="formatPrice" :parser="parsePrice" placeholder="最小价格" style="width: 45%" />
              <span class="separator">-</span>
              <el-input-number v-model="searchForm.maxPrice" :min="0.0" :step="50000" :formatter="formatPrice" :parser="parsePrice" placeholder="最大价格" style="width: 45%" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <div style="text-align: right; margin-top: 20px;">
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </div>
      </el-form>
    </el-card>

    <!-- 房产列表 -->
    <el-card class="realty-cards-container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else-if="realtyList.length === 0" class="empty-data">
        <el-empty description="暂无房产数据" />
      </div>
      <div v-else class="realty-grid">
        <el-card 
          v-for="item in realtyList" 
          :key="item.realtyCert" 
          class="realty-card"
          @click="viewDetails(item)"
        >
          <div class="realty-image">
            <img :src="item.images[0] ? item.images[0] : 'http://localhost:8089/i/2025/04/28/680f3098ac95a.png'" alt="房产图片" />
            <div class="realty-status" :class="getStatusClass(item.status)">{{ getStatusText(item.status) }}</div>
            <div class="is-new-house" v-if="item.isNewHouse">{{'新房'}}</div>
            <div class="is-not-new-house" v-if="!item.isNewHouse">{{'二手房'}}</div>
          </div>
          <div class="realty-info">
            <div class="realty-title">{{ generateAddress(item) }}</div>
            <div class="">{{ generateTitle(item) }}</div>
            <div class="realty-meta">
              <span class="house-type">{{ getHouseTypeText(item.houseType) }}</span>
              <span class="area">{{ item.area }}平米</span>
            </div>
            <div class="realty-price" v-if="userStore.hasOrganization(['investor', 'bank'])">
              <span>¥ {{ formatPriceText(item.price) }}</span>
            </div>
            <div class="realty-actions">
              <el-button type="primary" size="small" @click.stop="viewDetails(item)">详情</el-button>
              <el-button 
                v-if="userStore.hasOrganization('investor') && item.status === 'NORMAL'" 
                type="success" 
                size="small" 
                @click.stop="startTransaction(item)"
              >
                交易
              </el-button>
              <el-button 
                v-if="userStore.hasOrganization('government')" 
                type="warning" 
                size="small" 
                @click.stop="editRealty(item)"
              >
                编辑
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="searchForm.pageNumber"
          v-model:page-size="searchForm.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import {queryRealtyList} from "@/api/realty.ts";

const router = useRouter()
const userStore = useUserStore()

// 添加省份城市相关数据
const provinceList = ref([])
const cityList = ref([])

// 查询条件
const searchForm = reactive({
  realtyCert: '',
  realtyType: '',
  houseType: '',
  minPrice: null,
  maxPrice: null,
  minArea: null,
  maxArea: null,
  province: '',
  city: '',
  district: '',
  street: '',
  community: '',
  unit: '',
  floor: '',
  room: '',
  isNewHouse: null,
  pageSize: 10,
  pageNumber: 1,
  status: ''
})

// 重置查询条件
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    if (key === 'pageSize') {
      searchForm[key] = 10
    } else if (key === 'pageNumber') {
      searchForm[key] = 1
    }else if (key === 'isNewHouse') {
      searchForm[key] = null
    } else {
      searchForm[key] = key.startsWith('min') || key.startsWith('max') ? null : ''
    }
  })
  // 重置城市列表
  cityList.value = []
  handleSearch()
}

// 获取省份数据
const fetchProvinceList = async () => {
  provinceList.value = [
      { provinceCode: '11', provinceName: '北京市' },
      { provinceCode: '12', provinceName: '天津市' },
      { provinceCode: '13', provinceName: '河北省' },
      { provinceCode: '14', provinceName: '山西省' },
      { provinceCode: '15', provinceName: '内蒙古自治区' },
      { provinceCode: '21', provinceName: '辽宁省' },
      { provinceCode: '22', provinceName: '吉林省' },
      { provinceCode: '23', provinceName: '黑龙江省' },
      { provinceCode: '31', provinceName: '上海市' },
      { provinceCode: '32', provinceName: '江苏省' },
      { provinceCode: '33', provinceName: '浙江省' },
      { provinceCode: '34', provinceName: '安徽省' },
      { provinceCode: '35', provinceName: '福建省' },
      { provinceCode: '36', provinceName: '江西省' },
      { provinceCode: '37', provinceName: '山东省' },
      { provinceCode: '41', provinceName: '河南省' },
      { provinceCode: '42', provinceName: '湖北省' },
      { provinceCode: '43', provinceName: '湖南省' },
      { provinceCode: '44', provinceName: '广东省' },
      { provinceCode: '45', provinceName: '广西壮族自治区' },
      { provinceCode: '46', provinceName: '海南省' },
      { provinceCode: '50', provinceName: '重庆市' },
      { provinceCode: '51', provinceName: '四川省' },
      { provinceCode: '52', provinceName: '贵州省' },
      { provinceCode: '53', provinceName: '云南省' },
      { provinceCode: '54', provinceName: '西藏自治区' },
      { provinceCode: '61', provinceName: '陕西省' },
      { provinceCode: '62', provinceName: '甘肃省' },
      { provinceCode: '63', provinceName: '青海省' },
      { provinceCode: '64', provinceName: '宁夏回族自治区' },
      { provinceCode: '65', provinceName: '新疆维吾尔自治区' },
      { provinceCode: '71', provinceName: '台湾省' },
      { provinceCode: '81', provinceName: '香港特别行政区' },
      { provinceCode: '82', provinceName: '澳门特别行政区' }
    ]
}

// 处理省份变化
const handleProvinceChange = (provinceName) => {
  searchForm.city = ''
  searchForm.district = ''
  
  // 构建城市列表
  if (provinceName.includes('市')) {
    // 对于直辖市，城市就是省份名
    cityList.value = [provinceName]
    searchForm.city = provinceName
  } else {
    // 根据省份模拟城市列表（实际项目中应该从API获取）
    switch (provinceName) {
      case '河北省':
        cityList.value = ['石家庄市', '唐山市', '秦皇岛市', '邯郸市', '邢台市', '保定市', '张家口市', '承德市', '沧州市', '廊坊市', '衡水市']
        break
      case '山西省':
        cityList.value = ['太原市', '大同市', '阳泉市', '长治市', '晋城市', '朔州市', '晋中市', '运城市', '忻州市', '临汾市', '吕梁市']
        break
      case '内蒙古自治区':
        cityList.value = ['呼和浩特市', '包头市', '乌海市', '赤峰市', '通辽市', '鄂尔多斯市', '呼伦贝尔市', '巴彦淖尔市', '乌兰察布市', '兴安盟', '锡林郭勒盟', '阿拉善盟']
        break
      case '辽宁省':
        cityList.value = ['沈阳市', '大连市', '鞍山市', '抚顺市', '本溪市', '丹东市', '锦州市', '营口市', '阜新市', '辽阳市', '盘锦市', '铁岭市', '朝阳市', '葫芦岛市']
        break
      case '吉林省':
        cityList.value = ['长春市', '吉林市', '四平市', '辽源市', '通化市', '白山市', '松原市', '白城市', '延边朝鲜族自治州']
        break
      case '黑龙江省':
        cityList.value = ['哈尔滨市', '齐齐哈尔市', '牡丹江市', '佳木斯市', '大庆市', '绥化市', '鹤岗市', '鸡西市', '黑河市', '双鸭山市', '伊春市', '七台河市', '大兴安岭地区']
        break
      case '江苏省':
        cityList.value = ['南京市', '无锡市', '徐州市', '常州市', '苏州市', '南通市', '连云港市', '淮安市', '盐城市', '扬州市', '镇江市', '泰州市', '宿迁市']
        break
      case '浙江省':
        cityList.value = ['杭州市', '宁波市', '温州市', '嘉兴市', '湖州市', '绍兴市', '金华市', '衢州市', '舟山市', '台州市', '丽水市']
        break
      case '安徽省':
        cityList.value = ['合肥市', '芜湖市', '蚌埠市', '淮南市', '马鞍山市', '淮北市', '铜陵市', '安庆市', '黄山市', '滁州市', '阜阳市', '宿州市', '六安市', '亳州市', '池州市', '宣城市']
        break
      case '福建省':
        cityList.value = ['福州市', '厦门市', '莆田市', '三明市', '泉州市', '漳州市', '南平市', '龙岩市', '宁德市']
        break
      case '江西省':
        cityList.value = ['南昌市', '景德镇市', '萍乡市', '九江市', '新余市', '鹰潭市', '赣州市', '吉安市', '宜春市', '抚州市', '上饶市']
        break
      case '山东省':
        cityList.value = ['济南市', '青岛市', '淄博市', '枣庄市', '东营市', '烟台市', '潍坊市', '济宁市', '泰安市', '威海市', '日照市', '临沂市', '德州市', '聊城市', '滨州市', '菏泽市']
        break
      case '河南省':
        cityList.value = ['郑州市', '开封市', '洛阳市', '平顶山市', '安阳市', '鹤壁市', '新乡市', '焦作市', '濮阳市', '许昌市', '漯河市', '三门峡市', '南阳市', '商丘市', '信阳市', '周口市', '驻马店市', '济源市']
        break
      case '湖北省':
        cityList.value = ['武汉市', '黄石市', '十堰市', '宜昌市', '襄阳市', '鄂州市', '荆门市', '孝感市', '荆州市', '黄冈市', '咸宁市', '随州市', '恩施土家族苗族自治州', '仙桃市', '潜江市', '天门市', '神农架林区']
        break
      case '湖南省':
        cityList.value = ['长沙市', '株洲市', '湘潭市', '衡阳市', '邵阳市', '岳阳市', '常德市', '张家界市', '益阳市', '郴州市', '永州市', '怀化市', '娄底市', '湘西土家族苗族自治州']
        break
      case '广东省':
        cityList.value = ['广州市', '韶关市', '深圳市', '珠海市', '汕头市', '佛山市', '江门市', '湛江市', '茂名市', '肇庆市', '惠州市', '梅州市', '汕尾市', '河源市', '阳江市', '清远市', '东莞市', '中山市', '潮州市', '揭阳市', '云浮市']
        break
      case '广西壮族自治区':
        cityList.value = ['南宁市', '柳州市', '桂林市', '梧州市', '北海市', '防城港市', '钦州市', '贵港市', '玉林市', '百色市', '贺州市', '河池市', '来宾市', '崇左市']
        break
      case '海南省':
        cityList.value = ['海口市', '三亚市', '三沙市', '儋州市']
        break
      case '四川省':
        cityList.value = ['成都市', '自贡市', '攀枝花市', '泸州市', '德阳市', '绵阳市', '广元市', '遂宁市', '内江市', '乐山市', '南充市', '眉山市', '宜宾市', '广安市', '达州市', '雅安市', '巴中市', '资阳市', '阿坝藏族羌族自治州', '甘孜藏族自治州', '凉山彝族自治州']
        break
      case '贵州省':
        cityList.value = ['贵阳市', '六盘水市', '遵义市', '安顺市', '毕节市', '铜仁市', '黔西南布依族苗族自治州', '黔东南苗族侗族自治州', '黔南布依族苗族自治州']
        break
      case '云南省':
        cityList.value = ['昆明市', '曲靖市', '玉溪市', '保山市', '昭通市', '丽江市', '普洱市', '临沧市', '楚雄彝族自治州', '红河哈尼族彝族自治州', '文山壮族苗族自治州', '西双版纳傣族自治州', '大理白族自治州', '德宏傣族景颇族自治州', '怒江傈僳族自治州', '迪庆藏族自治州']
        break
      case '西藏自治区':
        cityList.value = ['拉萨市', '日喀则市', '昌都市', '林芝市', '山南市', '那曲市', '阿里地区']
        break
      case '陕西省':
        cityList.value = ['西安市', '铜川市', '宝鸡市', '咸阳市', '渭南市', '延安市', '汉中市', '榆林市', '安康市', '商洛市']
        break
      case '甘肃省':
        cityList.value = ['兰州市', '嘉峪关市', '金昌市', '白银市', '天水市', '武威市', '张掖市', '平凉市', '酒泉市', '庆阳市', '定西市', '陇南市', '临夏回族自治州', '甘南藏族自治州']
        break
      case '青海省':
        cityList.value = ['西宁市', '海东市', '海北藏族自治州', '黄南藏族自治州', '海南藏族自治州', '果洛藏族自治州', '玉树藏族自治州', '海西蒙古族藏族自治州']
        break
      case '宁夏回族自治区':
        cityList.value = ['银川市', '石嘴山市', '吴忠市', '固原市', '中卫市']
        break
      case '新疆维吾尔自治区':
        cityList.value = ['乌鲁木齐市', '克拉玛依市', '吐鲁番市', '哈密市', '昌吉回族自治州', '博尔塔拉蒙古自治州', '巴音郭楞蒙古自治州', '阿克苏地区', '克孜勒苏柯尔克孜自治州', '喀什地区', '和田地区', '伊犁哈萨克自治州', '塔城地区', '阿勒泰地区', '石河子市', '阿拉尔市', '图木舒克市', '五家渠市', '北屯市', '铁门关市', '双河市', '可克达拉市', '昆玉市', '胡杨河市', '新星市']
        break
      case '台湾省':
        cityList.value = ['台北市', '高雄市', '台中市', '台南市', '新北市', '桃园市', '基隆市', '新竹市', '嘉义市']
        break
      case '香港特别行政区':
        cityList.value = ['中西区', '东区', '南区', '湾仔区', '九龙城区', '观塘区', '深水埗区', '黄大仙区', '油尖旺区', '离岛区', '葵青区', '北区', '西贡区', '沙田区', '大埔区', '荃湾区', '屯门区', '元朗区']
        break
      case '澳门特别行政区':
        cityList.value = ['花地玛堂区', '圣安多尼堂区', '大堂区', '望德堂区', '风顺堂区', '嘉模堂区', '路凼填海区', '圣方济各堂区']
        break
      default:
        cityList.value = []
    }
  }
}

// 处理城市变化
const handleCityChange = (cityName) => {
  searchForm.district = ''
}

// 格式化价格显示
const formatPrice = (value) => {
  if (value === null) return ''
  return `¥ ${value}`
}

// 解析价格
const parsePrice = (value) => {
  if (value === '') return null
  return value.replace(/[^\d]/g, '')
}

// 价格文本格式化
const formatPriceText = (price) => {
  if (!price) return '暂无报价'
  return (price / 10000).toFixed(2) + '万'
}

// 获取随机图片
const getRandomImage = (id) => {
  // 使用房产证号作为种子生成一个稳定的随机数，这样同一个房产始终显示相同的图片
  const seed = id.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  const imageId = (seed % 100) + 1 // 限制在1-100之间
  return `https://picsum.photos/id/${imageId}/300/200`
}

// 获取状态样式类
const getStatusClass = (status) => {
  const statusMap = {
    'NORMAL': 'status-normal',
    'IN_SALE': 'status-pending',
    'MORTGAGED': 'status-mortgaged',
    'PENDING_SALE': 'status-pendingSale',
    'FROZEN': 'status-frozen'
  }
  return statusMap[status] || ''
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'IN_SALE': '交易中',
    'IN_MORTGAGE': '已抵押',
    'PENDING_SALE': '挂牌中',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || '未知状态'
}

// 获取户型文本
const getHouseTypeText = (houseType) => {
  const houseTypeMap = {
    'single': '一室',
    'double': '两室',
    'triple': '三室',
    'multiple': '四室及以上'
  }
  return houseTypeMap[houseType] || houseType
}

// 生成房产标题
const generateTitle = (item) => {
  if (!item) return '未知房产'
  
  let title = ''
  
  // 房产类型
  const typeMap = {
    'HOUSE': '住宅',
    'SHOP': '商铺',
    'OFFICE': '办公',
    'INDUSTRIAL': '工业',
    'OTHER': '其他'
  }
  
  title += typeMap[item.realtyType] || '未知类型'
  
  // 户型信息
  if (item.realtyType === 'HOUSE') {
    title += ' - ' + getHouseTypeText(item.houseType)
  }
  
  // 小区/位置信息
  if (item.community) {
    title += ' - ' + item.community
  }
  
  return title
}

// 生成房产地址
const generateAddress = (item) => {
  if (!item) return '地址不详'
  
  const parts = []
  if (item.province) parts.push(item.province)
  if (item.city && item.city !== item.province) parts.push(item.city)
  if (item.district && item.district !== item.city) parts.push(item.district)
  if (item.street) parts.push(item.street)
  if (item.community) parts.push(item.community)
  if (item.unit) parts.push(item.unit + '单元')
  if (item.floor) parts.push(item.floor + '层')
  if (item.room) parts.push(item.room + '室')
  
  return parts.length > 0 ? parts.join(' ') : '地址不详'
}

// 房产列表数据
const realtyList = ref([])
const loading = ref(true)
const total = ref(0)

// 获取房产列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await queryRealtyList({
      realtyCert: searchForm.realtyCert,
      realtyType: searchForm.realtyType,
      houseType: searchForm.houseType,
      minPrice: searchForm.minPrice,
      maxPrice: searchForm.maxPrice,
      minArea: searchForm.minArea,
      maxArea: searchForm.maxArea,
      province: searchForm.province,
      city: searchForm.city,
      district: searchForm.district,
      street: searchForm.street,
      community: searchForm.community,
      unit: searchForm.unit,
      isNewHouse: searchForm.isNewHouse,
      floor: searchForm.floor,
      room: searchForm.room,
      pageSize: searchForm.pageSize,
      pageNumber: searchForm.pageNumber,
      status: searchForm.status,
    })

    // 适配新的API返回格式
    if (response && response.realtyList) {
      realtyList.value = response.realtyList
      total.value = response.total || 0
    } else {
      realtyList.value = []
      total.value = 0
    }
  } catch (error) {
    ElMessage.error('获取房产列表失败，请稍后再试')
    realtyList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = () => {
  searchForm.pageNumber = 1
  fetchData()
}

// 改变每页显示数量
const handleSizeChange = (size) => {
  searchForm.pageSize = size
  fetchData()
}

// 改变页码
const handleCurrentChange = (page) => {
  searchForm.pageNumber = page
  fetchData()
}

// 查看详情
const viewDetails = (item) => {
  router.push(`/realty/${item.realtyCertHash}`)
}

// 开始交易
const startTransaction = (item) => {
  router.push({
    path: '/transaction/create',
    query: { realtyCert: item.realtyCert }
  })
}

// 编辑房产
const editRealty = (item) => {
  router.push(`/realty/${item.id}`)
}

// 初始加载
onMounted(() => {
  fetchData()
  // 获取省份数据
  fetchProvinceList()
})
</script>

<style scoped>
.realty-list {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 22px;
}

.filter-container {
  margin-bottom: 20px;
}

.separator {
  margin: 0 5px;
}

.realty-cards-container {
  margin-top: 20px;
}

.loading-container {
  padding: 20px 0;
}

.empty-data {
  padding: 40px 0;
  text-align: center;
}

.realty-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.realty-card {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;
}

.realty-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.realty-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.realty-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.realty-status {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.is-new-house {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: rgba(0, 100, 0, 0.7);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.is-not-new-house {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: rgba(200, 100, 0, 0.7);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-normal {
  background-color: rgba(25, 190, 107, 0.7);
}

.status-pending {
  background-color: #E6A23C;
}

.status-mortgaged {
  background-color: #409EFF;
}

.status-frozen {
  background-color: #F56C6C;
}

.realty-info {
  padding: 15px;
}

.realty-title {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.realty-address {
  color: #606266;
  font-size: 14px;
  margin-bottom: 10px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.realty-meta {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  color: #909399;
  font-size: 13px;
}

.realty-price {
  font-size: 18px;
  color: #F56C6C;
  font-weight: bold;
  margin: 10px 0;
}

.realty-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 15px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
