<template>
  <div class="teaching">
    <div class="page-header">
      <h2>课程管理</h2>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        创建课程
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="courses"
      style="width: 100%; margin-top: 20px"
    >
      <el-table-column prop="name" label="课程名称" min-width="200">
        <template #default="{ row }">
          <el-button text type="primary" @click="showEditDialog(row)">
            {{ row.name }}
          </el-button>
        </template>
      </el-table-column>
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
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            @click="showStudentList(row)"
          >
            学生名单
          </el-button>
          <el-button
            type="warning"
            size="small"
            @click="toggleCourseStatus(row)"
          >
            {{ row.status === 'open' ? '关闭' : '开放' }}
          </el-button>
          <el-button
            type="danger"
            size="small"
            @click="handleDelete(row)"
          >
            删除
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

    <!-- 创建/编辑课程对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑课程' : '创建课程'"
      width="50%"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="课程描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="学分" prop="credits">
          <el-input-number
            v-model="form.credits"
            :min="1"
            :max="4"
            :precision="1"
          />
        </el-form-item>
        <el-form-item label="容量" prop="capacity">
          <el-input-number
            v-model="form.capacity"
            :min="1"
            :max="200"
            :precision="0"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 学生名单对话框 -->
    <el-dialog
      v-model="studentDialogVisible"
      title="学生名单"
      width="60%"
    >
      <el-table
        v-loading="loadingStudents"
        :data="students"
        style="width: 100%"
      >
        <el-table-column prop="student.username" label="学生姓名" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '在读' : '已退课' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="成绩" width="200">
          <template #default="{ row }">
            <el-input-number
              v-if="row.status === 'active'"
              v-model="row.grade"
              :min="0"
              :max="100"
              :precision="1"
              @change="(value) => handleGradeChange(row, value)"
            />
            <span v-else>{{ row.grade || '-' }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const loading = ref(false)
const loadingStudents = ref(false)
const submitting = ref(false)
const courses = ref([])
const students = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const dialogVisible = ref(false)
const studentDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  name: '',
  description: '',
  credits: 1,
  capacity: 50
})

const rules = {
  name: [
    { required: true, message: '请输入课程名称', trigger: 'blur' },
    { min: 2, message: '课程名称至少2个字符', trigger: 'blur' }
  ],
  credits: [
    { required: true, message: '请输入学分', trigger: 'blur' },
    { type: 'number', min: 1, max: 4, message: '学分必须在1-4之间', trigger: 'blur' }
  ],
  capacity: [
    { required: true, message: '请输入容量', trigger: 'blur' },
    { type: 'number', min: 1, max: 200, message: '容量必须在1-200之间', trigger: 'blur' }
  ]
}

// 获取课程列表
async function fetchCourses() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value
    }

    const response = await axios.get('/api/v1/teaching/courses', { params })
    courses.value = response.data.courses
    total.value = response.data.total
  } catch (error) {
    console.error('Failed to fetch courses:', error)
    ElMessage.error('获取课程列表失败')
  } finally {
    loading.value = false
  }
}

// 显示创建对话框
function showCreateDialog() {
  isEdit.value = false
  Object.assign(form, {
    name: '',
    description: '',
    credits: 1,
    capacity: 50
  })
  dialogVisible.value = true
}

// 显示编辑对话框
function showEditDialog(course) {
  isEdit.value = true
  Object.assign(form, {
    name: course.name,
    description: course.description,
    credits: course.credits,
    capacity: course.capacity
  })
  dialogVisible.value = true
}

// 提交表单
async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value) {
      await axios.put(`/api/v1/courses/${course.id}`, form)
      ElMessage.success('课程更新成功')
    } else {
      await axios.post('/api/v1/courses', form)
      ElMessage.success('课程创建成功')
    }

    dialogVisible.value = false
    fetchCourses()
  } catch (error) {
    console.error('Submit failed:', error)
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 切换课程状态
async function toggleCourseStatus(course) {
  try {
    const newStatus = course.status === 'open' ? 'closed' : 'open'
    await axios.put(`/api/v1/courses/${course.id}/status`, {
      status: newStatus
    })
    course.status = newStatus
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('Toggle status failed:', error)
    ElMessage.error('状态更新失败')
  }
}

// 删除课程
async function handleDelete(course) {
  try {
    await ElMessageBox.confirm(
      '删除课程后无法恢复，是否继续？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await axios.delete(`/api/v1/courses/${course.id}`)
    ElMessage.success('课程删除成功')
    fetchCourses()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 显示学生名单
async function showStudentList(course) {
  loadingStudents.value = true
  studentDialogVisible.value = true
  try {
    const response = await axios.get(`/api/v1/courses/${course.id}/students`)
    students.value = response.data.students
  } catch (error) {
    console.error('Failed to fetch students:', error)
    ElMessage.error('获取学生名单失败')
  } finally {
    loadingStudents.value = false
  }
}

// 更新学生成绩
async function handleGradeChange(student, grade) {
  try {
    await axios.put(`/api/v1/enrollments/${student.id}/grade`, {
      grade
    })
    ElMessage.success('成绩更新成功')
  } catch (error) {
    console.error('Update grade failed:', error)
    ElMessage.error('成绩更新失败')
    // 恢复原值
    student.grade = student.original_grade
  }
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
.teaching {
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 