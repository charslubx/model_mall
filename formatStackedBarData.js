/**
 * 将嵌套对象数据转换为 ECharts 柱状图 series 配置
 * @param {Object} data - 原始数据对象
 * @param {Array} colors - 颜色数组 [color0, color1, color2, color3]
 * @param {Number} type - 图表类型：1=堆叠柱状图, 2=分组柱状图
 * @param {Object} options - 可选配置
 * @returns {Array} ECharts series 配置数组
 */
function formatStackedBarData(data, colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666'], type = 1, options = {}) {
  const {
    // 定义数据类别的顺序（外层 key 的顺序）
    categoryOrder = [
      'Actual',
      'IF_ATM_Down_Time',
      'IPG_XRB',
      'IF_FSM_DMO_XRB',
      'IF_XRB',
      'IF_NON_XRB',
      'IF_EDISPOSE_REWORK'
    ],
    // 定义堆叠层的配置
    stackLayers = [
      { name: 'Assy', field: 'Assy', colorIndex: 0 },
      { name: 'Test_ATPO', field: 'test_atpo', colorIndex: 1 },
      { name: 'Test_FPO', field: 'test_fpo', colorIndex: 2 },
      { name: 'Finish', field: 'finish', colorIndex: 3 }
    ],
    stackName = 'st',
    showLabel = true
  } = options;

  // 根据类型决定是否需要倒序
  const layers = type === 1 ? [...stackLayers].reverse() : stackLayers;

  // 构建 series 数组
  const series = layers.map(layer => {
    // 为每个类别提取对应的值
    const seriesData = categoryOrder.map(category => {
      const categoryData = data[category];
      if (!categoryData) return null;
      
      const value = categoryData[layer.field];
      return value !== undefined ? value : null;
    });

    // 基础配置
    const seriesItem = {
      name: layer.name,
      type: 'bar',
      itemStyle: { color: colors[layer.colorIndex] },
      data: seriesData
    };

    // type=1: 堆叠柱状图配置
    if (type === 1) {
      seriesItem.stack = stackName;
      seriesItem.emphasis = { focus: 'series' };
      
      if (showLabel) {
        seriesItem.label = {
          show: true,
          position: 'inside',
          color: '#000',
          backgroundColor: 'rgba(255,255,255,0.65)',
          borderRadius: 2,
          padding: [1, 3],
          textStyle: {
            fontSize: 10,
          },
          formatter: p => {
            const v = p.data;
            if (v === null || v === undefined) return '';
            return v >= 100 ? (v >= 1000 ? (v / 1000).toFixed(1) + 'k' : v) : '';
          }
        };
      }
    }
    // type=2: 分组柱状图配置
    else if (type === 2) {
      if (showLabel) {
        seriesItem.label = {
          show: true,
          position: 'top',
          formatter: '{c}'
        };
      }
    }

    return seriesItem;
  });

  return series;
}

/**
 * 获取 x 轴类别数据
 * @param {Object} data - 原始数据对象
 * @param {Array} categoryOrder - 类别顺序
 * @returns {Array} x 轴类别数组
 */
function getXAxisCategories(data, categoryOrder = null) {
  if (categoryOrder) {
    return categoryOrder;
  }
  return Object.keys(data);
}

// ============= 测试示例 =============

const testData = {
  "Actual": {
    "Assy": 8253,
    "test_atpo": 8237,
    "test_fpo": 10677,
    "finish": 10913
  },
  "IF_ATM_Down_Time": {
    "Assy": 118,
    "test_atpo": 56,
    "test_fpo": 1,
    "finish": 73
  },
  "IPG_XRB": {
    "Assy": 0,
    "test_atpo": 0,
    "test_fpo": 0,
    "finish": 0
  },
  "IF_FSM_DMO_XRB": {
    "Assy": 8,
    "test_atpo": 2,
    "test_fpo": 1,
    "finish": 8
  },
  "IF_XRB": {
    "Assy": 84,
    "test_atpo": 254,
    "test_fpo": 47,
    "finish": 634
  },
  "IF_NON_XRB": {
    "Assy": 5449,
    "test_atpo": 757,
    "test_fpo": 400,
    "finish": 1415
  },
  "IF_EDISPOSE_REWORK": {
    "test_atpo": 3
  }
};

// 定义颜色（根据你的需求调整）
const colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666'];

// ========== 测试 Type 1: 堆叠柱状图 ==========
console.log('========== Type 1: 堆叠柱状图 ==========');
const seriesType1 = formatStackedBarData(testData, colors, 1);

console.log('生成的 series 配置：');
console.log(JSON.stringify(seriesType1, null, 2));

console.log('\n数据验证 (Type 1):');
seriesType1.forEach(s => {
  console.log(`  ${s.name}: [${s.data.map(v => v === null ? 'null' : v).join(', ')}]`);
});

// ========== 测试 Type 2: 分组柱状图 ==========
console.log('\n\n========== Type 2: 分组柱状图 ==========');
const seriesType2 = formatStackedBarData(testData, colors, 2);

console.log('生成的 series 配置：');
console.log(JSON.stringify(seriesType2, null, 2));

console.log('\n数据验证 (Type 2):');
seriesType2.forEach(s => {
  console.log(`  ${s.name}: [${s.data.map(v => v === null ? 'null' : v).join(', ')}]`);
});

// 生成 x 轴类别
const xAxisData = getXAxisCategories(testData);

console.log('\n\nx 轴类别：');
console.log(JSON.stringify(xAxisData, null, 2));

console.log('\n\n========== 完整的 ECharts 配置示例 (Type 1) ==========');
const completeOptionType1 = {
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    }
  },
  legend: {
    data: ['Finish', 'Test_FPO', 'Test_ATPO', 'Assy']
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: xAxisData,
    axisLabel: {
      rotate: 45,
      interval: 0
    }
  },
  yAxis: {
    type: 'value'
  },
  series: seriesType1
};

console.log(JSON.stringify(completeOptionType1, null, 2));

console.log('\n\n========== 完整的 ECharts 配置示例 (Type 2) ==========');
const completeOptionType2 = {
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    }
  },
  legend: {
    data: ['Assy', 'Test_ATPO', 'Test_FPO', 'Finish']
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: xAxisData,
    axisLabel: {
      rotate: 45,
      interval: 0
    }
  },
  yAxis: {
    type: 'value'
  },
  series: seriesType2
};

console.log(JSON.stringify(completeOptionType2, null, 2));

// 导出函数
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    formatStackedBarData,
    getXAxisCategories
  };
}
