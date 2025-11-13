# ECharts markLine 超出数据范围的解决方案

## 问题描述
当 markLine 的 y 值（17）大于数据实际最大值（例如只有 5）时，ECharts 会自动缩放 y 轴，导致 markLine 无法显示。

## 解决方案

### 方案 1：使用函数动态设置 y 轴最大值（推荐 ⭐）

这是最灵活的方案，可以自动适应不同的数据情况。

```javascript
optionDSDP1 = {
  ...this.optionTmp,
  ...this.optionXY,
  yAxis: {
    type: 'value',
    // 动态计算 y 轴最大值，确保包含 markLine
    max: function(value) {
      const markLineValue = 17;  // 你的 markLine 值
      // 取数据最大值和 markLine 值中的较大者，并留出 10% 的空间
      return Math.ceil(Math.max(value.max, markLineValue) * 1.1);
    },
    min: 0  // 可选：设置最小值为 0
  },
  series: [
    {
      name: 'Lead Time',
      type: 'bar',
      data: [
        {value: 36.1, itemStyle: {color: '#d65e5e'}},
        {value: 36.1, itemStyle: {color: '#d65e5e'}},
        {value: 36.1, itemStyle: {color: '#d65e5e'}},
        {value: 27.1, itemStyle: {color: '#d65e5e'}},
        {value: 25.0, itemStyle: {color: '#d65e5e'}},
        {value: 11.7, itemStyle: {color: '#d65e5e'}},
        {value: 9.0, itemStyle: {color: '#d65e5e'}}
      ],
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
```

### 方案 2：固定 y 轴范围（最简单）

如果你的数据范围相对固定，可以直接设置固定值。

```javascript
optionDSDP1 = {
  ...this.optionTmp,
  ...this.optionXY,
  yAxis: {
    type: 'value',
    min: 0,
    max: 40  // 固定最大值，确保 markLine (17) 可见
  },
  series: [
    // ... 其他配置不变
  ]
};
```

### 方案 3：计算后设置（适合需要更多控制的场景）

```javascript
// 在设置 option 之前计算 y 轴范围
function getYAxisRange(data, markLineValue) {
  const values = data.map(item => item.value || 0);
  const dataMax = Math.max(...values);
  const effectiveMax = Math.max(dataMax, markLineValue);
  
  return {
    min: 0,
    max: Math.ceil(effectiveMax * 1.1)  // 留出 10% 空间
  };
}

const yAxisRange = getYAxisRange(yourData, 17);

optionDSDP1 = {
  ...this.optionTmp,
  ...this.optionXY,
  yAxis: {
    type: 'value',
    min: yAxisRange.min,
    max: yAxisRange.max
  },
  series: [
    // ... 其他配置不变
  ]
};
```

## 推荐选择

- **如果数据动态变化**：使用 **方案 1**（函数方式）
- **如果数据范围固定**：使用 **方案 2**（固定值）
- **如果需要复杂逻辑**：使用 **方案 3**（计算方式）

## 额外优化

如果你想让 markLine 总是显示在合适的位置，还可以：

```javascript
yAxis: {
  type: 'value',
  min: 0,
  max: function(value) {
    const markLineValue = 17;
    const dataMax = value.max;
    
    // 如果数据最大值远小于 markLine（比如小于 markLine 的 50%）
    if (dataMax < markLineValue * 0.5) {
      // 使用 markLine 的 1.2 倍作为最大值
      return Math.ceil(markLineValue * 1.2);
    } else {
      // 否则使用数据最大值的 1.1 倍
      return Math.ceil(Math.max(dataMax, markLineValue) * 1.1);
    }
  }
}
```

这样可以确保在任何情况下，markLine 都能清晰可见！
