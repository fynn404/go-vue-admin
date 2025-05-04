<!-- 成绩统计组件 -->
<template>
  <div class="grade-stats">
    <!-- 基础统计信息 -->
    <el-card class="stats-card">
      <template #header>
        <div class="card-header">
          <span>{{ title }}</span>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :span="6" v-for="(item, index) in basicStats" :key="index">
          <div class="stat-item">
            <div class="stat-value">{{ item.value }}</div>
            <div class="stat-label">{{ item.label }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 成绩分布图表 -->
    <el-card class="chart-card">
      <template #header>
        <div class="card-header">
          <span>成绩分布</span>
        </div>
      </template>
      <div class="chart-container">
        <div ref="distributionChart" style="height: 300px"></div>
      </div>
    </el-card>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { onMounted, ref, watch } from 'vue'

export default {
  name: 'GradeStats',
  props: {
    title: {
      type: String,
      default: '成绩统计'
    },
    stats: {
      type: Object,
      required: true
    }
  },
  setup(props) {
    const distributionChart = ref(null)
    let chart = null

    // 处理基础统计数据
    const basicStats = ref([])
    
    const updateBasicStats = () => {
      if (props.stats.courseStats) {
        basicStats.value = [
          { label: '总人数', value: props.stats.totalStudents },
          { label: '平均分', value: props.stats.averageScore?.toFixed(1) },
          { label: '最高分', value: props.stats.highestScore?.toFixed(1) },
          { label: '及格率', value: props.stats.passRate?.toFixed(1) + '%' }
        ]
      } else {
        // 学生统计
        basicStats.value = [
          { label: '课程总数', value: props.stats.courseCount },
          { label: '已通过课程', value: props.stats.passedCourses },
          { label: '总学分', value: props.stats.totalCredits },
          { label: 'GPA', value: props.stats.gpa?.toFixed(2) }
        ]
      }
    }

    // 初始化图表
    const initChart = () => {
      if (!distributionChart.value) return
      
      chart = echarts.init(distributionChart.value)
      updateChart()
    }

    // 更新图表数据
    const updateChart = () => {
      if (!chart || !props.stats.distribution) return

      const option = {
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: props.stats.distribution.map(item => item.range)
        },
        yAxis: [
          {
            type: 'value',
            name: '人数',
            position: 'left'
          },
          {
            type: 'value',
            name: '百分比',
            position: 'right',
            axisLabel: {
              formatter: '{value}%'
            }
          }
        ],
        series: [
          {
            name: '人数',
            type: 'bar',
            data: props.stats.distribution.map(item => item.count)
          },
          {
            name: '百分比',
            type: 'line',
            yAxisIndex: 1,
            data: props.stats.distribution.map(item => item.percent)
          }
        ]
      }

      chart.setOption(option)
    }

    // 监听数据变化
    watch(() => props.stats, () => {
      updateBasicStats()
      updateChart()
    }, { deep: true })

    // 组件挂载时初始化
    onMounted(() => {
      updateBasicStats()
      initChart()

      // 响应式调整
      window.addEventListener('resize', () => {
        chart?.resize()
      })
    })

    return {
      distributionChart,
      basicStats
    }
  }
}
</script>

<style scoped>
.grade-stats {
  padding: 20px;
}

.stats-card {
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}

.stat-label {
  margin-top: 10px;
  color: #666;
}

.chart-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 