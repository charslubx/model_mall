/**
 * 将对象数据按照指定顺序格式化为数组
 * @param {Object} data - 原始数据对象
 * @param {Object} options - 可选配置
 * @param {Array} options.order - 属性顺序数组（可选）
 * @param {Object} options.colorMap - 颜色映射对象（可选）
 * @returns {Array} 格式化后的数组
 */
function formatAttributeData(data, options = {}) {
  // 默认顺序
  const defaultOrder = [
    'Actual',
    'IF_ATM_Down_Time',
    'IPG_XRB',
    'IF_FSM_DMO_XRB',
    'IF_XRB',
    'IF_NON_XRB',
    'IF_EDISPOSE_REWORK'
  ];

  // 默认颜色映射
  const defaultColorMap = {
    'Actual': '#d65e5e',
    'IF_ATM_Down_Time': '#f39c12',
    'IPG_XRB': '#3498db',
    'IF_FSM_DMO_XRB': '#9b59b6',
    'IF_XRB': '#1abc9c',
    'IF_NON_XRB': '#e74c3c',
    'IF_EDISPOSE_REWORK': '#95a5a6'
  };

  // 使用传入的配置或默认配置
  const order = options.order || defaultOrder;
  const colorMap = options.colorMap || defaultColorMap;

  // 按照顺序构建结果数组
  const result = [];
  
  order.forEach(attr => {
    // 如果数据中存在该属性
    if (data.hasOwnProperty(attr)) {
      result.push({
        attr: attr,
        value: data[attr],
        itemStyle: {
          color: colorMap[attr] || '#d65e5e' // 如果没有指定颜色，使用默认颜色
        }
      });
    }
  });

  return result;
}

// 示例使用
const testData = {
  "Actual": "5.04",
  "IF_ATM_Down_Time": "4.96",
  "IPG_XRB": "4.96",
  "IF_FSM_DMO_XRB": "4.96",
  "IF_XRB": "4.86",
  "IF_NON_XRB": "4.10",
  "IF_EDISPOSE_REWORK": "2.17"
};

console.log('默认配置输出：');
console.log(JSON.stringify(formatAttributeData(testData), null, 2));

console.log('\n使用自定义颜色：');
const customColors = {
  'Actual': '#ff0000',
  'IF_ATM_Down_Time': '#00ff00',
  'IPG_XRB': '#0000ff',
  'IF_FSM_DMO_XRB': '#ffff00',
  'IF_XRB': '#ff00ff',
  'IF_NON_XRB': '#00ffff',
  'IF_EDISPOSE_REWORK': '#888888'
};
console.log(JSON.stringify(formatAttributeData(testData, { colorMap: customColors }), null, 2));

console.log('\n使用统一颜色：');
const singleColor = {
  'Actual': '#d65e5e',
  'IF_ATM_Down_Time': '#d65e5e',
  'IPG_XRB': '#d65e5e',
  'IF_FSM_DMO_XRB': '#d65e5e',
  'IF_XRB': '#d65e5e',
  'IF_NON_XRB': '#d65e5e',
  'IF_EDISPOSE_REWORK': '#d65e5e'
};
console.log(JSON.stringify(formatAttributeData(testData, { colorMap: singleColor }), null, 2));

// 导出函数（用于Node.js模块）
if (typeof module !== 'undefined' && module.exports) {
  module.exports = formatAttributeData;
}

// 导出函数（用于ES6模块）
// export default formatAttributeData;
