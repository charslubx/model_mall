/**
 * 将嵌套对象数据转换为 ECharts 堆叠柱状图 series 配置
 * @param {Object} data - 原始数据对象
 * @param {Array} colors - 颜色数组 [color0, color1, color2, color3]
 * @param {Object} options - 可选配置
 * @returns {Array} ECharts series 配置数组
 */
function formatStackedBarData(data, colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666'], options = {}) {
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
    // 定义堆叠层的配置（从下到上）
    stackLayers = [
      { name: 'Assy', field: 'Assy', colorIndex: 0 },
      { name: 'Test_ATPO', field: 'test_atpo', colorIndex: 1 },
      { name: 'Test_FPO', field: 'test_fpo', colorIndex: 2 },
      { name: 'Finish', field: 'finish', colorIndex: 3 }
    ],
    stackName = 'st',
    showLabel = true
  } = options;

  // 构建 series 数组（注意：需要倒序，因为 ECharts 堆叠是从数组开始往后叠加的）
  const series = stackLayers.reverse().map(layer => {
    // 为每个类别提取对应的值
    const seriesData = categoryOrder.map(category => {
      const categoryData = data[category];
      if (!categoryData) return null;
      
      const value = categoryData[layer.field];
      return value !== undefined ? value : null;
    });

    return {
      name: layer.name,
      type: 'bar',
      stack: stackName,
      emphasis: { focus: 'series' },
      itemStyle: { color: colors[layer.colorIndex] },
      data: seriesData,
      label: showLabel ? {
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
      } : undefined
    };
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

// 生成 series 配置
const series = formatStackedBarData(testData, colors);

// 生成 x 轴类别
const xAxisData = getXAxisCategories(testData);

console.log('生成的 series 配置：');
console.log(JSON.stringify(series, null, 2));

console.log('\n\nx 轴类别：');
console.log(JSON.stringify(xAxisData, null, 2));

console.log('\n\n完整的 ECharts 配置示例：');
const completeOption = {
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
  series: series
};

console.log(JSON.stringify(completeOption, null, 2));

// 验证数据
console.log('\n\n=== 数据验证 ===');
series.forEach(s => {
  console.log(`\n${s.name}:`);
  console.log(`  数据: [${s.data.join(', ')}]`);
  console.log(`  颜色: ${s.itemStyle.color}`);
});

// 导出函数
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    formatStackedBarData,
    getXAxisCategories
  };
}
