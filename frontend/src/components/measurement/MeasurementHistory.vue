<template>
  <div class="measurement-history">
    <el-card>
      <template #header>
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium">测量历史记录</h3>
          <div>
            <el-button type="primary" size="small" @click="refresh">刷新</el-button>
            <el-button type="success" size="small" @click="navigateToCreate">新建测量</el-button>
          </div>
        </div>
      </template>
      
      <!-- 统计卡片 -->
      <div v-if="stats" class="stats-card mb-6 p-4 bg-gray-50 rounded-lg">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="stat-item p-3 bg-white rounded shadow-sm">
            <div class="text-sm text-gray-500">手型分类</div>
            <div class="text-xl font-bold">{{ handSizeLabel }}</div>
          </div>
          <div class="stat-item p-3 bg-white rounded shadow-sm">
            <div class="text-sm text-gray-500">平均手掌宽度</div>
            <div class="text-xl font-bold">{{ stats.averagePalm }} mm</div>
          </div>
          <div class="stat-item p-3 bg-white rounded shadow-sm">
            <div class="text-sm text-gray-500">平均手指长度</div>
            <div class="text-xl font-bold">{{ stats.averageLength }} mm</div>
          </div>
        </div>
        <div class="text-right mt-2 text-sm text-gray-500">
          测量次数: {{ stats.measurementCount }} | 最后测量: {{ formatDate(stats.lastMeasuredAt) }}
        </div>
      </div>
      
      <!-- 筛选 -->
      <div class="filters mb-4">
        <el-form :model="filters" inline>
          <el-form-item label="排序">
            <el-select v-model="filters.sortBy" size="small" @change="fetchData">
              <el-option label="创建时间" value="createdAt"></el-option>
              <el-option label="手掌宽度" value="palm"></el-option>
              <el-option label="手指长度" value="length"></el-option>
              <el-option label="测量质量" value="quality"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="顺序">
            <el-select v-model="filters.sortOrder" size="small" @change="fetchData">
              <el-option label="降序" value="desc"></el-option>
              <el-option label="升序" value="asc"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="日期范围">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small"
              @change="handleDateChange"
            ></el-date-picker>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 数据表格 -->
      <el-table 
        :data="measurements" 
        style="width: 100%" 
        border 
        v-loading="loading"
        :empty-text="emptyText"
      >
        <el-table-column prop="createdAt" label="测量日期" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="palm" label="手掌宽度" width="120">
          <template #default="scope">
            {{ formatMeasurement(scope.row.palm, scope.row.unit) }}
          </template>
        </el-table-column>
        <el-table-column prop="length" label="手指长度" width="120">
          <template #default="scope">
            {{ formatMeasurement(scope.row.length, scope.row.unit) }}
          </template>
        </el-table-column>
        <el-table-column prop="quality" label="测量质量" width="120">
          <template #default="scope">
            <el-progress 
              :percentage="scope.row.quality ? scope.row.quality.score : 0"
              :status="getQualityStatus(scope.row.quality ? scope.row.quality.score : 0)"
            ></el-progress>
          </template>
        </el-table-column>
        <el-table-column label="校准状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.quality && scope.row.quality.factors.calibration ? 'success' : 'warning'">
              {{ scope.row.quality && scope.row.quality.factors.calibration ? '已校准' : '未校准' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="viewDetails(scope.row)">详情</el-button>
            <el-button size="small" type="danger" @click="deleteMeasurement(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
      
      <!-- 详情对话框 -->
      <el-dialog v-model="detailsVisible" title="测量详情" width="500px">
        <div v-if="selectedMeasurement" class="measurement-details">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <div class="text-sm text-gray-500">手掌宽度</div>
              <div class="text-lg font-bold">{{ formatMeasurement(selectedMeasurement.palm, selectedMeasurement.unit) }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">手指长度</div>
              <div class="text-lg font-bold">{{ formatMeasurement(selectedMeasurement.length, selectedMeasurement.unit) }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">测量质量</div>
              <div class="text-lg font-bold">{{ selectedMeasurement.quality ? selectedMeasurement.quality.score : 0 }}%</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">校准状态</div>
              <div class="text-lg font-bold">{{ selectedMeasurement.quality && selectedMeasurement.quality.factors.calibration ? '已校准' : '未校准' }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">创建时间</div>
              <div class="text-lg font-bold">{{ formatDateTime(selectedMeasurement.createdAt) }}</div>
            </div>
            <div>
              <div class="text-sm text-gray-500">更新时间</div>
              <div class="text-lg font-bold">{{ formatDateTime(selectedMeasurement.updatedAt) }}</div>
            </div>
          </div>
          
          <div class="mt-4 p-3 bg-gray-50 rounded-lg">
            <div class="text-sm font-medium mb-2">测量质量详情</div>
            <template v-if="selectedMeasurement.quality">
              <div class="grid grid-cols-2 gap-2">
                <div>
                  <div class="text-xs text-gray-500">校准</div>
                  <div>{{ selectedMeasurement.quality.factors.calibration ? '是' : '否' }}</div>
                </div>
                <div>
                  <div class="text-xs text-gray-500">稳定性</div>
                  <div>{{ (selectedMeasurement.quality.factors.stability * 100).toFixed(0) }}%</div>
                </div>
                <div>
                  <div class="text-xs text-gray-500">一致性</div>
                  <div>{{ (selectedMeasurement.quality.factors.consistency * 100).toFixed(0) }}%</div>
                </div>
              </div>
            </template>
            <div v-else class="text-sm text-gray-400">无质量数据</div>
          </div>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { format } from 'date-fns'
import { useMeasurement } from '@/composables/useMeasurement'

// 路由
const router = useRouter()

// 状态
const loading = ref(false)
const measurements = ref<any[]>([]) // 这里应该使用正确的类型
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const stats = ref<any>(null) // 用户统计数据
const dateRange = ref<[Date, Date] | null>(null)
const detailsVisible = ref(false)
const selectedMeasurement = ref<any>(null)

// 空状态文本
const emptyText = computed(() => {
  return loading.value ? '加载中...' : '暂无测量记录'
})

// 过滤器
const filters = reactive({
  sortBy: 'createdAt',
  sortOrder: 'desc',
  startDate: '',
  endDate: ''
})

// 手型分类标签
const handSizeLabel = computed(() => {
  if (!stats.value) return '未知'
  
  switch (stats.value.handSize) {
    case 'small': return '小型手'
    case 'medium': return '中型手'
    case 'large': return '大型手'
    default: return '未知'
  }
})

// 初始化
onMounted(() => {
  fetchData()
  fetchStats()
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    // 调用实际的API
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      sortBy: filters.sortBy,
      sortOrder: filters.sortOrder,
      startDate: filters.startDate,
      endDate: filters.endDate
    }
    
    const { fetchMeasurements } = useMeasurement()
    const response = await fetchMeasurements(params)
    
    if (response) {
      measurements.value = response.measurements
      total.value = response.total
    } else {
      // 如果API调用失败，使用模拟数据
      console.warn('API调用失败，使用模拟数据')
      measurements.value = [
        {
          id: '1',
          palm: 85.5,
          length: 75.2,
          unit: 'mm',
          quality: {
            score: 85,
            factors: {
              calibration: true,
              stability: 0.8,
              consistency: 0.9
            }
          },
          createdAt: new Date('2023-01-01T10:00:00'),
          updatedAt: new Date('2023-01-01T10:00:00')
        },
        {
          id: '2',
          palm: 90.2,
          length: 78.8,
          unit: 'mm',
          quality: {
            score: 70,
            factors: {
              calibration: false,
              stability: 0.7,
              consistency: 0.8
            }
          },
          createdAt: new Date('2023-01-05T15:30:00'),
          updatedAt: new Date('2023-01-05T15:30:00')
        }
      ]
      total.value = 2
    }
  } catch (error) {
    console.error('获取测量记录失败', error)
    ElMessage.error('获取记录失败，请刷新重试')
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    // 调用实际的API
    const { fetchMeasurementStats } = useMeasurement()
    const response = await fetchMeasurementStats()
    
    if (response) {
      stats.value = response
    } else {
      // 如果API调用失败，使用模拟数据
      console.warn('API调用失败，使用模拟数据')
      stats.value = {
        averagePalm: 87.8,
        averageLength: 77.0,
        handSize: 'medium',
        measurementCount: 2,
        lastMeasuredAt: new Date('2023-01-05T15:30:00')
      }
    }
  } catch (error) {
    console.error('获取统计数据失败', error)
    // 静默失败，不显示错误提示
  }
}

// 刷新数据
const refresh = () => {
  fetchData()
  fetchStats()
}

// 导航到创建页面
const navigateToCreate = () => {
  router.push('/measurements/create')
}

// 处理日期变化
const handleDateChange = (val: [Date, Date] | null) => {
  if (val) {
    filters.startDate = format(val[0], 'yyyy-MM-dd')
    filters.endDate = format(val[1], 'yyyy-MM-dd')
  } else {
    filters.startDate = ''
    filters.endDate = ''
  }
  fetchData()
}

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchData()
}

// 处理页码变化
const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchData()
}

// 查看详情
const viewDetails = (row: any) => {
  selectedMeasurement.value = row
  detailsVisible.value = true
}

// 删除测量记录
const deleteMeasurement = (row: any) => {
  ElMessageBox.confirm(
    '确定要删除这条测量记录吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 调用API删除记录
      const { removeMeasurement } = useMeasurement()
      const success = await removeMeasurement(row.id)
      
      if (success) {
        ElMessage.success('删除成功')
        fetchData() // 刷新数据
        fetchStats() // 刷新统计数据
      } else {
        ElMessage.error('删除失败，请重试')
      }
    } catch (error) {
      console.error('删除失败', error)
      ElMessage.error('删除失败，请重试')
    }
  }).catch(() => {
    // 取消删除，不做任何操作
  })
}

// 格式化日期
const formatDate = (date: string | Date) => {
  if (!date) return '—'
  return format(new Date(date), 'yyyy-MM-dd')
}

// 格式化日期时间
const formatDateTime = (date: string | Date) => {
  if (!date) return '—'
  return format(new Date(date), 'yyyy-MM-dd HH:mm:ss')
}

// 格式化测量值
const formatMeasurement = (value: number, unit: string) => {
  return `${value.toFixed(1)} ${unit}`
}

// 获取质量状态
const getQualityStatus = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'exception'
}
</script>

<style scoped>
.stat-item {
  transition: all 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}
</style>