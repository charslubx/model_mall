# formatStackedBarData 使用说明

## 方法签名

```javascript
formatStackedBarData(data, colors, type, options)
```

## 参数说明

- **data** (Object): 原始数据对象
- **colors** (Array): 颜色数组 `[color0, color1, color2, color3]`
- **type** (Number): 图表类型
  - `1` = 堆叠柱状图（Stacked Bar Chart）
  - `2` = 分组柱状图（Grouped Bar Chart）
- **options** (Object, 可选): 额外配置选项

## 使用示例

### Type 1: 堆叠柱状图

```javascript
const data = {
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
  // ... 其他数据
};

const colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666'];

// 生成堆叠柱状图 series
const series = formatStackedBarData(data, colors, 1);

// 使用在 ECharts option 中
const option = {
  ...this.optionTmp,
  ...this.optionXY,
  series: series
};
```

**Type 1 输出特点：**
- 有 `stack: 'st'` 属性（柱子堆叠显示）
- 有 `emphasis: { focus: 'series' }` 
- label 位置在 `inside`（柱子内部）
- label 有背景色和样式
- label 数值格式化（≥1000 显示为 "k"）
- series 顺序：Finish → Test_FPO → Test_ATPO → Assy（从上到下）

**输出示例：**
```javascript
[
  {
    name: 'Finish',
    type: 'bar',
    stack: 'st',
    emphasis: { focus: 'series' },
    itemStyle: { color: '#ee6666' },
    data: [10913, 73, 0, 8, 634, 1415, null],
    label: {
      show: true,
      position: 'inside',
      color: '#000',
      backgroundColor: 'rgba(255,255,255,0.65)',
      borderRadius: 2,
      padding: [1, 3],
      textStyle: { fontSize: 10 },
      formatter: p => { /* ... */ }
    }
  },
  // ... 其他 series
]
```

---

### Type 2: 分组柱状图

```javascript
const data = { /* 同上 */ };
const colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666'];

// 生成分组柱状图 series
const series = formatStackedBarData(data, colors, 2);

// 使用在 ECharts option 中
const option = {
  ...this.optionTmp,
  ...this.optionXY,
  series: series
};
```

**Type 2 输出特点：**
- 没有 `stack` 属性（柱子并排显示）
- label 位置在 `top`（柱子顶部）
- label 简单格式化 `{c}`
- series 顺序：Assy → Test_ATPO → Test_FPO → Finish（从左到右）

**输出示例：**
```javascript
[
  {
    name: 'Assy',
    type: 'bar',
    itemStyle: { color: '#5470c6' },
    data: [8253, 118, 0, 8, 84, 5449, null],
    label: {
      show: true,
      position: 'top',
      formatter: '{c}'
    }
  },
  {
    name: 'Test_ATPO',
    type: 'bar',
    itemStyle: { color: '#91cc75' },
    data: [8237, 56, 0, 2, 254, 757, 3],
    label: {
      show: true,
      position: 'top',
      formatter: '{c}'
    }
  },
  {
    name: 'Test_FPO',
    type: 'bar',
    itemStyle: { color: '#fac858' },
    data: [10677, 1, 0, 1, 47, 400, null],
    label: {
      show: true,
      position: 'top',
      formatter: '{c}'
    }
  },
  {
    name: 'Finish',
    type: 'bar',
    itemStyle: { color: '#ee6666' },
    data: [10913, 73, 0, 8, 634, 1415, null],
    label: {
      show: true,
      position: 'top',
      formatter: '{c}'
    }
  }
]
```

---

## 数据映射规则

### x 轴类别顺序（固定）
```javascript
[
  'Actual',
  'IF_ATM_Down_Time',
  'IPG_XRB',
  'IF_FSM_DMO_XRB',
  'IF_XRB',
  'IF_NON_XRB',
  'IF_EDISPOSE_REWORK'
]
```

### series 层级与字段映射
```javascript
[
  { name: 'Assy',      field: 'Assy',      colorIndex: 0 },
  { name: 'Test_ATPO', field: 'test_atpo', colorIndex: 1 },
  { name: 'Test_FPO',  field: 'test_fpo',  colorIndex: 2 },
  { name: 'Finish',    field: 'finish',    colorIndex: 3 }
]
```

### 空值处理
- 如果某个字段不存在或值为 `undefined`，则设置为 `null`
- 例如：`IF_EDISPOSE_REWORK` 没有 `Assy` 字段，则该位置为 `null`

---

## 在 Vue 组件中使用

```javascript
export default {
  data() {
    return {
      colors: ['#5470c6', '#91cc75', '#fac858', '#ee6666'],
      rawData: { /* 你的数据 */ }
    };
  },
  methods: {
    // 复制 formatStackedBarData 方法到这里
    formatStackedBarData(data, colors, type, options = {}) {
      // ... 方法实现
    },
    
    // 更新图表
    updateChart(type) {
      const series = this.formatStackedBarData(this.rawData, this.colors, type);
      
      this.optionDSDP1 = {
        ...this.optionTmp,
        ...this.optionXY,
        series: series
      };
    }
  },
  mounted() {
    // 初始化为堆叠柱状图
    this.updateChart(1);
  }
}
```

---

## 完整配置示例

```javascript
const option = {
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow'
    }
  },
  legend: {
    // Type 1: ['Finish', 'Test_FPO', 'Test_ATPO', 'Assy']
    // Type 2: ['Assy', 'Test_ATPO', 'Test_FPO', 'Finish']
    data: type === 1 
      ? ['Finish', 'Test_FPO', 'Test_ATPO', 'Assy']
      : ['Assy', 'Test_ATPO', 'Test_FPO', 'Finish']
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: [
      'Actual', 'IF_ATM_Down_Time', 'IPG_XRB', 
      'IF_FSM_DMO_XRB', 'IF_XRB', 'IF_NON_XRB', 
      'IF_EDISPOSE_REWORK'
    ],
    axisLabel: {
      rotate: 45,
      interval: 0
    }
  },
  yAxis: {
    type: 'value'
  },
  series: formatStackedBarData(data, colors, type)
};
```

---

## 两种类型的对比

| 特性 | Type 1 (堆叠) | Type 2 (分组) |
|------|--------------|--------------|
| 柱子排列 | 垂直堆叠 | 水平并排 |
| stack 属性 | 有 (`'st'`) | 无 |
| label 位置 | `inside` | `top` |
| label 样式 | 有背景色、padding | 无背景色 |
| label 格式化 | 智能格式化 (k) | 简单显示 `{c}` |
| emphasis | 有 | 无 |
| series 顺序 | 倒序（上→下） | 正序（左→右） |
| 适用场景 | 显示总量构成 | 对比各项数值 |
