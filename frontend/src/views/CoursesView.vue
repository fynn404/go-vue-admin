<template>
  <div class="courses">
    <div class="page-header">
      <h2>选课中心</h2>
      <el-input
        v-model="searchQuery"
        placeholder="搜索课程"
        prefix-icon="Search"
        clearable
        @input="handleSearch"
        style="width: 300px"
      />
    </div>

    <el-card class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="课程状态">
          <el-select v-model="filters.status" placeholder="全部" clearable>
            <el-option label="开放" value="open" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="学分">
          <el-select v-model="filters.credits" placeholder="全部" clearable>
            <el-option label="1学分" value="1" />
            <el-option label="2学分" value="2" />
            <el-option label="3学分" value="3" />
            <el-option label="4学分" value="4" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">筛选</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-table
      v-loading="loading"
      :data="courses"
      style="width: 100%; margin-top: 20px"
    >
      <el-table-column prop="name" label="课程名称" min-width="200">
        <template #default="{ row }">
          <el-button text type="primary" @click="showCourseDetails(row)">
            {{ row.name }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="teacher.username" label="授课教师" />
      <el-table-column prop="credits" label="学分" width="100" />
      <el-table-column prop="capacity" label="容量" width="150">
        <template #default="{ row }">
          {{ row.current_enrolled }}/{{ row.capacity }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'open' ? 'success' : 'info'">
            {{ row.status === 'open' ? '开放' : '已关闭' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            :disabled="!row.is_available"
            @click="handleEnroll(row)"
          >
            选课
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 课程详情对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="课程详情"
      width="50%"
      destroy-on-close
    >
      <template v-if="selectedCourse">
        <div class="course-details">
          <h3>{{ selectedCourse.name }}</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="授课教师">
              {{ selectedCourse.teacher.username }}
            </el-descriptions-item>
            <el-descriptions-item label="学分">
              {{ selectedCourse.credits }}
            </el-descriptions-item>
            <el-descriptions-item label="容量">
              {{ selectedCourse.current_enrolled }}/{{ selectedCourse.capacity }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="selectedCourse.status === 'open' ? 'success' : 'info'">
                {{ selectedCourse.status === 'open' ? '开放' : '已关闭' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
          <div class="course-description">
            <h4>课程描述</h4>
            <p>{{ selectedCourse.description || '暂无描述' }}</p>
          </div>
        </div>
      </template>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
          <el-button
            type="primary"
            :disabled="!selectedCourse?.is_available"
            @click="handleEnroll(selectedCourse)"
          >
            选课
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const loading = ref(false)
const courses = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const dialogVisible = ref(false)
const selectedCourse = ref(null)

const filters = reactive({
  status: '',
  credits: ''
})

// 获取课程列表
async function fetchCourses() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      search: searchQuery.value,
      status: filters.status,
      credits: filters.credits
    }

    const response = await axios.get('/api/v1/courses', { params })
    courses.value = response.data.courses
    total.value = response.data.total
  } catch (error) {
    console.error('Failed to fetch courses:', error)
    ElMessage.error('获取课程列表失败')
  } finally {
    loading.value = false
  }
}

// 处理选课
async function handleEnroll(course) {
  try {
    await ElMessageBox.confirm(
      `确定要选修课程 "${course.name}" 吗？`,
      '选课确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    const response = await axios.post(`/api/v1/courses/${course.id}/enroll`)
    ElMessage.success('选课成功')
    dialogVisible.value = false
    fetchCourses() // 刷新课程列表
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Enrollment failed:', error)
      ElMessage.error(error.response?.data?.error || '选课失败')
    }
  }
}

// 显示课程详情
function showCourseDetails(course) {
  selectedCourse.value = course
  dialogVisible.value = true
}

// 处理搜索
const handleSearch = _.debounce(() => {
  currentPage.value = 1
  fetchCourses()
}, 300)

// 处理筛选
function handleFilter() {
  currentPage.value = 1
  fetchCourses()
}

// 重置筛选
function resetFilters() {
  filters.status = ''
  filters.credits = ''
  currentPage.value = 1
  fetchCourses()
}

// 处理分页
function handleSizeChange(val) {
  pageSize.value = val
  fetchCourses()
}

function handleCurrentChange(val) {
  currentPage.value = val
  fetchCourses()
}

onMounted(() => {
  fetchCourses()
})
</script>

<style scoped>
.courses {
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
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.course-details {
  padding: 20px;
}

.course-details h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #409EFF;
}

.course-description {
  margin-top: 20px;
}

.course-description h4 {
  margin-bottom: 10px;
  color: #606266;
}

.course-description p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 