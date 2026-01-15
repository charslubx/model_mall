package types

// AnalyzeAnomalyRequest 前端点击图表点位后上报的上下文
// 说明：
// - xValue/yValue 对应被点击的数据点
// - seriesValues/xValues 建议携带“同一条曲线”的完整序列或窗口序列，用于后端做统计对比/构造AI提示词
type AnalyzeAnomalyRequest struct {
	Chart       string            `json:"chart,omitempty"`       // 图表标识/标题（可选）
	Metric      string            `json:"metric,omitempty"`      // 指标名（可选）
	SeriesName  string            `json:"seriesName,omitempty"`  // 系列名（可选）
	XName       string            `json:"xName,omitempty"`       // x轴字段名（可选）
	XValue      string            `json:"xValue,omitempty"`      // 点击点的x（类目/日期字符串等）
	YName       string            `json:"yName,omitempty"`       // y轴字段名（可选）
	YValue      float64           `json:"yValue"`                // 点击点的y值
	DataIndex   int               `json:"dataIndex,omitempty"`   // 点击点在序列中的索引（可选）
	SeriesValues []float64        `json:"seriesValues,omitempty"`// 同系列数值序列（可选）
	XValues     []string          `json:"xValues,omitempty"`     // 同系列x序列（可选）
	TimeRange   string            `json:"timeRange,omitempty"`   // 查询窗口（可选，比如 7days/30days）
	Filters     map[string]string `json:"filters,omitempty"`     // 业务过滤条件（可选）
	Extra       map[string]any    `json:"extra,omitempty"`       // 其他扩展字段（可选）
}

type AnalyzeAnomalyResponse struct {
	Summary         string   `json:"summary"`
	Causes          []string `json:"causes,omitempty"`
	Actions         []string `json:"actions,omitempty"`
	Confidence      float64  `json:"confidence,omitempty"`
	UsedModel       bool     `json:"usedModel,omitempty"`
	ModelProvider   string   `json:"modelProvider,omitempty"`
}

