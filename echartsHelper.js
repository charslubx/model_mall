/**
 * 计算 ECharts y 轴的合适范围，确保 markLine 可见
 * @param {Array} data - 图表数据数组
 * @param {Number} markLineValue - markLine 的 y 轴值
 * @param {Object} options - 可选配置
 * @returns {Object} 包含 min 和 max 的对象
 */
function calculateYAxisRange(data, markLineValue, options = {}) {
  const {
    padding = 0.1,        // 上下留白比例（10%）
    minBuffer = 2         // 最小缓冲值
  } = options;

  // 提取所有数值
  const values = data.map(item => {
    if (typeof item === 'number') return item;
    if (item && typeof item.value === 'number') return item.value;
    return 0;
  });

  // 计算数据的最大值和最小值
  const dataMax = Math.max(...values, 0);
  const dataMin = Math.min(...values, 0);

  // 确保最大值至少包含 markLine 的值
  const effectiveMax = Math.max(dataMax, markLineValue);
  
  // 添加缓冲区
  const range = effectiveMax - dataMin;
  const bufferValue = Math.max(range * padding, minBuffer);

  return {
    min: Math.floor(dataMin - bufferValue),
    max: Math.ceil(effectiveMax + bufferValue)
  };
}

/**
 * 示例 1：基本用法
 */
const exampleData1 = [
  {value: 5, itemStyle: {color: '#d65e5e'}},
  {value: 4.5, itemStyle: {color: '#d65e5e'}},
  {value: 3.8, itemStyle: {color: '#d65e5e'}},
  {value: 2.7, itemStyle: {color: '#d65e5e'}},
  {value: 2.0, itemStyle: {color: '#d65e5e'}},
];

const markLineValue1 = 17;
const yAxisRange1 = calculateYAxisRange(exampleData1, markLineValue1);

console.log('示例 1 - 数据最大值小于 markLine:');
console.log('数据最大值: 5, markLine 值: 17');
console.log('计算的 y 轴范围:', yAxisRange1);
console.log('');

// ECharts 配置示例 1
const optionExample1 = {
  xAxis: {
    type: 'category',
    data: ['Item1', 'Item2', 'Item3', 'Item4', 'Item5']
  },
  yAxis: {
    type: 'value',
    min: yAxisRange1.min,
    max: yAxisRange1.max,
    // 或者使用函数动态设置
    // max: function(value) {
    //   return Math.max(value.max, markLineValue1) * 1.1;
    // }
  },
  series: [{
    name: 'Lead Time',
    type: 'bar',
    data: exampleData1,
    barWidth: '70%',
    label: {
      show: true,
      position: 'inside',
      fontSize: 12,
      color: '#ffffff',
      formatter: '{c}'
    },
    markLine: {
      silent: true,
      symbol: 'none',
      data: [{
        yAxis: markLineValue1,
        lineStyle: {
          type: 'dashed',
          color: '#72b4dd',
          width: 2
        },
        label: {
          show: true,
          position: 'insideStartTop',
          formatter: markLineValue1.toFixed(1),
          color: '#72b4dd',
          fontSize: 14,
          fontWeight: 'bold',
          distance: 5
        }
      }]
    }
  }]
};

console.log('ECharts 配置示例（简化版）:');
console.log(JSON.stringify({
  yAxis: {
    type: 'value',
    min: yAxisRange1.min,
    max: yAxisRange1.max
  }
}, null, 2));
console.log('');

/**
 * 示例 2：数据最大值大于 markLine
 */
const exampleData2 = [
  {value: 36.1, itemStyle: {color: '#d65e5e'}},
  {value: 36.1, itemStyle: {color: '#d65e5e'}},
  {value: 36.1, itemStyle: {color: '#d65e5e'}},
  {value: 27.1, itemStyle: {color: '#d65e5e'}},
  {value: 25.0, itemStyle: {color: '#d65e5e'}},
  {value: 11.7, itemStyle: {color: '#d65e5e'}},
  {value: 9.0, itemStyle: {color: '#d65e5e'}}
];

const markLineValue2 = 17;
const yAxisRange2 = calculateYAxisRange(exampleData2, markLineValue2);

console.log('示例 2 - 数据最大值大于 markLine:');
console.log('数据最大值: 36.1, markLine 值: 17');
console.log('计算的 y 轴范围:', yAxisRange2);
console.log('');

/**
 * 方案 2：使用 ECharts 内置的 max 函数（推荐）
 */
console.log('=== 推荐方案：直接在 yAxis 中使用函数 ===');
console.log(`
yAxis: {
  type: 'value',
  max: function(value) {
    // value.max 是数据中的最大值
    // 确保 y 轴最大值至少包含 markLine 的值，并留出 10% 的空间
    const markLineValue = 17;
    return Math.ceil(Math.max(value.max, markLineValue) * 1.1);
  },
  min: function(value) {
    // 可选：设置最小值
    return Math.floor(value.min * 0.9);
  }
}
`);

/**
 * 方案 3：更简单的固定范围方案
 */
console.log('=== 简单方案：固定 y 轴范围 ===');
console.log(`
yAxis: {
  type: 'value',
  min: 0,
  max: 40  // 固定最大值，确保 markLine 可见
}
`);

/**
 * 完整的解决方案示例
 */
const completeOption = {
  title: {
    text: 'Lead Time Analysis'
  },
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
    data: ['Actual', 'IF_ATM_Down_Time', 'IPG_XRB', 'IF_FSM_DMO_XRB', 'IF_XRB', 'IF_NON_XRB', 'IF_EDISPOSE_REWORK'],
    axisLabel: {
      rotate: 45,
      interval: 0
    }
  },
  yAxis: {
    type: 'value',
    name: 'Days',
    // 方法 1: 使用计算好的值
    // min: yAxisRange1.min,
    // max: yAxisRange1.max,
    
    // 方法 2: 使用函数（推荐）
    max: function(value) {
      const markLineValue = 17;
      return Math.ceil(Math.max(value.max, markLineValue) * 1.1);
    },
    min: 0
  },
  series: [
    {
      name: 'Lead Time',
      type: 'bar',
      data: exampleData1,  // 使用小数据测试
      barWidth: '70%',
      label: {
        show: true,
        position: 'inside',
        fontSize: 12,
        color: '#ffffff',
        formatter: '{c}',
        overflow: 'none',
        distance: 0
      },
      itemStyle: {
        borderRadius: 0
      },
      markLine: {
        silent: true,
        symbol: 'none',
        data: [
          {
            yAxis: 17,
            lineStyle: {
              type: 'dashed',
              color: '#72b4dd',
              width: 2
            },
            label: {
              show: true,
              position: 'insideStartTop',
              formatter: '17.0',
              color: '#72b4dd',
              fontSize: 14,
              fontWeight: 'bold',
              distance: 5
            }
          }
        ]
      }
    }
  ]
};

console.log('=== 完整配置已生成 ===');
console.log('配置对象变量名: completeOption');

// 导出函数
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    calculateYAxisRange,
    completeOption
  };
}
