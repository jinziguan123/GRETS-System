<template>
  <div class="realty-create-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <h3>创建新房产</h3>
        </div>
      </template>
      <el-form 
        :model="realtyForm" 
        :rules="rules" 
        ref="realtyFormRef" 
        label-width="120px"
        label-position="right"
        status-icon
      >
        <!-- 基础信息 -->
        <h4>基础信息</h4>
        <el-form-item label="不动产证号" prop="realtyCert">
          <el-input v-model="realtyForm.realtyCert" placeholder="请输入不动产证号"></el-input>
        </el-form-item>
        
        <el-form-item label="房产类型" prop="realtyType">
          <el-select v-model="realtyForm.realtyType" placeholder="请选择房产类型" style="width: 100%">
            <el-option label="住宅" value="HOUSE"></el-option>
            <el-option label="商铺" value="SHOP"></el-option>
            <el-option label="办公" value="OFFICE"></el-option>
            <el-option label="工业" value="INDUSTRIAL"></el-option>
            <el-option label="其他" value="OTHER"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="户型" prop="houseType">
          <el-select v-model="realtyForm.houseType" placeholder="请选择户型" style="width: 100%">
            <el-option label="一室" value="single"></el-option>
            <el-option label="两室" value="double"></el-option>
            <el-option label="三室" value="triple"></el-option>
            <el-option label="四室及以上" value="multiple"></el-option>
          </el-select>
        </el-form-item>

        <el-row>
          <el-col :span="12">
            <el-form-item label="面积" prop="area">
              <el-input-number v-model="realtyForm.area" :min="0" :precision="2" :step="1" style="width: 100%" placeholder="请输入面积（平方米）"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格" prop="price">
              <el-input-number v-model="realtyForm.price" :min="0" :precision="2" :step="10000" style="width: 100%" placeholder="请输入价格（元）"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>

<!--        <el-form-item label="当前所有者" prop="currentOwnerCitizenID">-->
<!--          <el-input v-model="realtyForm.currentOwnerCitizenID" placeholder="请输入当前所有者身份证号"></el-input>-->
<!--        </el-form-item>-->

        <el-form-item label="状态" prop="status">
          <el-select v-model="realtyForm.status" placeholder="请选择房产状态" style="width: 100%">
            <el-option label="正常" value="NORMAL"></el-option>
            <el-option label="挂牌" value="PENDING_SALE"></el-option>
            <el-option label="已抵押" value="IN_MORTGAGE"></el-option>
            <el-option label="已冻结" value="FROZEN"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="关联合同" prop="contractID" v-if="realtyForm.status === 'PENDING_SALE'">
          <el-select v-model="realtyForm.contractUUID" placeholder="请选择关联合同" style="width: 100%" filterable>
            <el-option 
              v-for="contract in contractList" 
              :key="contract.contractUUID" 
              :label="contract.title" 
              :value="contract.contractUUID" 
            ></el-option>
          </el-select>
          <div class="form-tip">新房设置为"正常"状态时，必须选择一份合同</div>
        </el-form-item>
        
        <!-- 地址信息 -->
        <h4>地址信息</h4>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="省份" prop="province">
              <el-select 
                v-model="realtyForm.province" 
                placeholder="请选择省份" 
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
            <el-form-item label="城市" prop="city">
              <el-select 
                v-model="realtyForm.city" 
                placeholder="请选择城市" 
                style="width: 100%"
                :disabled="!realtyForm.province"
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
            <el-form-item label="区/县" prop="district">
              <el-input 
                v-model="realtyForm.district" 
                placeholder="请输入区/县"
                :disabled="!realtyForm.city"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="街道" prop="street">
              <el-input 
                v-model="realtyForm.street" 
                placeholder="请输入街道"
                :disabled="!realtyForm.district"
              ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="小区" prop="community">
              <el-input 
                v-model="realtyForm.community" 
                placeholder="请输入小区名称"
                :disabled="!realtyForm.street"
              ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="单元" prop="unit">
              <el-input 
                v-model="realtyForm.unit" 
                placeholder="请输入单元"
                :disabled="!realtyForm.community"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="楼层" prop="floor">
              <el-input 
                v-model="realtyForm.floor" 
                placeholder="请输入楼层"
                :disabled="!realtyForm.unit"
              ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房号" prop="room">
              <el-input 
                v-model="realtyForm.room" 
                placeholder="请输入房号"
                :disabled="!realtyForm.floor"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 附加信息 -->
        <h4>附加信息</h4>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="realtyForm.description" 
            type="textarea" 
            rows="3" 
            placeholder="请输入房产描述"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="房产图片" prop="images">
          <el-upload
            action="http://localhost:8080/api/v1/picture/upload"
            list-type="picture-card"
            :limit="5"
            :file-list="fileList"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :on-success="handleFileSuccess"
          >
            <el-icon><Plus /></el-icon>
            <template #tip>
              <div class="el-upload__tip">
                只能上传 jpg/png 文件，且不超过 5MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="loading">创建</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { createRealty } from "@/api/realty.js"
import { useUserStore } from '@/stores/user'
import { queryContractList } from '@/api/contract'
import { getProvinces } from '@/api/region'

const router = useRouter()
const userStore = useUserStore()
const realtyFormRef = ref(null)
const loading = ref(false)
const fileList = ref([])
const contractList = ref([])
const provinceList = ref([])
const cityList = ref([])

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
  realtyForm.city = ''
  realtyForm.district = ''
  realtyForm.street = ''
  realtyForm.community = ''
  realtyForm.unit = ''
  realtyForm.floor = ''
  realtyForm.room = ''
  
  // 构建城市列表
  if (provinceName.includes('市')) {
    // 对于直辖市，城市就是省份名
    cityList.value = [provinceName]
    realtyForm.city = provinceName
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

// 获取合同列表
const fetchContracts = async () => {
  try {
    const response = await queryContractList({
      status: 'NORMAL',
      pageSize: 100,
      pageNumber: 1,
      creatorCitizenID: userStore.user.citizenID,
    })
    contractList.value = response.contracts.filter(e => !e.transactionUUID) || []
  } catch (error) {
    console.error('Failed to fetch contracts:', error)
    ElMessage.error('获取合同列表失败')
  }
}

// 表单数据
const realtyForm = reactive({
  realtyCert: '',
  realtyType: '',
  houseType: '',
  price: 0,
  area: 0,
  currentOwnerCitizenID: userStore.user.citizenID, // 政府默认账户
  currentOwnerOrganization: 'government', // 默认为政府
  status: 'NORMAL',
  province: '',
  city: '',
  district: '',
  street: '',
  community: '',
  unit: '',
  floor: '',
  room: '',
  description: '',
  // 这些将作为计算属性或默认值
  address: '',
  images: [],
  previousOwnersCitizenIDList: [],
  contractUUID: '',
  isNewHouse: true,
})

// 计算完整地址
const updateAddress = () => {
  // 根据省份类型决定地址格式
  if (realtyForm.province.includes('市')) {
    // 直辖市格式
    realtyForm.address = `${realtyForm.province}${realtyForm.district}区${realtyForm.street}${realtyForm.community}${realtyForm.unit}单元${realtyForm.floor}楼${realtyForm.room}`
  } else {
    // 普通省份格式
    realtyForm.address = `${realtyForm.province}${realtyForm.city}${realtyForm.district}区${realtyForm.street}${realtyForm.community}${realtyForm.unit}单元${realtyForm.floor}楼${realtyForm.room}`
  }
}

// 监听地址相关字段变化
const addressFields = ['province', 'city', 'district', 'street', 'community', 'unit', 'floor', 'room']
addressFields.forEach(field => {
  watch(() => realtyForm[field], () => {
    updateAddress()
  })
})

// 监听状态变化，如果状态为NORMAL，则需要选择合同
watch(() => realtyForm.status, (newVal) => {
  if (newVal === 'PENDING_SALE') {
    // 修改规则，添加合同ID验证
    rules.contractUUID = [
      { required: true, message: '正常状态的新房必须选择关联合同', trigger: 'change' }
    ]
  } else {
    // 移除合同ID验证
    delete rules.contractUUID
  }
})

// 文件上传相关方法
const handleFileChange = (file) => {
  console.log('文件变化:', file)
  // 实际环境中，这里会上传文件并获取URL
  // 这里只是模拟
}

const handleFileRemove = (file) => {
  console.log('移除文件:', file)
}

const handleFileSuccess = (response) => {
  console.log('上传成功:', response)
  realtyForm.images.push(response.data.url)
}

// 表单验证规则
const rules = reactive({
  realtyCert: [
    { required: true, message: '请输入不动产证号', trigger: 'blur' },
    { min: 5, max: 50, message: '长度应在 5 到 50 个字符之间', trigger: 'blur' }
  ],
  realtyType: [
    { required: true, message: '请选择房产类型', trigger: 'change' }
  ],
  houseType: [
    { required: true, message: '请选择户型', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格必须大于等于0', trigger: 'blur' }
  ],
  area: [
    { required: true, message: '请输入面积', trigger: 'blur' },
    { type: 'number', min: 0, message: '面积必须大于等于0', trigger: 'blur' }
  ],
  currentOwnerCitizenID: [
    { required: true, message: '请输入当前所有者身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择房产状态', trigger: 'change' }
  ],
  province: [
    { required: true, message: '请选择省份', trigger: 'change' }
  ],
  city: [
    { required: true, message: '请选择城市', trigger: 'change' }
  ],
  district: [
    { required: true, message: '请输入区/县', trigger: 'blur' }
  ],
  street: [
    { required: true, message: '请输入街道', trigger: 'blur' }
  ],
  community: [
    { required: true, message: '请输入小区', trigger: 'blur' }
  ],
  unit: [
    { required: true, message: '请输入单元', trigger: 'blur' }
  ],
  floor: [
    { required: true, message: '请输入楼层', trigger: 'blur' }
  ],
  room: [
    { required: true, message: '请输入房号', trigger: 'blur' }
  ]
})

// 提交表单
const submitForm = async () => {
  if (!realtyFormRef.value) return
  
  await realtyFormRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning('请完善表单信息')
      return
    }
    
    // 检查状态为NORMAL时是否选择了合同
    if (realtyForm.status === 'PENDING_SALE' && !realtyForm.contractUUID) {
      ElMessage.warning('正常状态的新房必须选择关联合同')
      return
    }
    
    loading.value = true
    
    try {
      // 准备提交数据
      const requestData = {
        ...realtyForm,
        relContractUUID: realtyForm.contractUUID || '',
      }
      
      await createRealty(requestData)
      ElMessage.success('房产创建成功')
      router.push('/realty')
    } catch (error) {
      ElMessage.error('创建房产失败:', error.response.data.message)
    } finally {
      loading.value = false
    }
  })
}

// 重置表单
const resetForm = () => {
  if (realtyFormRef.value) {
    realtyFormRef.value.resetFields()
    fileList.value = []
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 组件初始化
onMounted(() => {
  // 获取合同列表
  fetchContracts()
  
  // 获取省份数据
  fetchProvinceList()
})
</script>

<style scoped>
.realty-create-container {
  max-width: 900px;
  margin: 20px auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.box-card {
  margin-bottom: 20px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

/* 卡片内部每个小标题的样式 */
h4 {
  margin-top: 20px;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
  color: #303133;
}
</style>
