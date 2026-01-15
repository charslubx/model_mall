package ai

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnalyzeAnomalyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点击图表点位后进行异常原因分析
func NewAnalyzeAnomalyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalyzeAnomalyLogic {
	return &AnalyzeAnomalyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AnalyzeAnomalyLogic) AnalyzeAnomaly(req *types.AnalyzeAnomalyRequest) (*types.AnalyzeAnomalyResponse, error) {
	// 需要登录（沿用现有 JWTMiddleware 在 ctx 中写入 userId）
	if _, ok := l.ctx.Value("userId").(int64); !ok {
		return nil, fmt.Errorf("未授权访问")
	}

	// 1) 优先调用模型服务（若模型服务实现了该端点）
	if l.svcCtx.ModelServiceClient != nil {
		modelResp, err := l.svcCtx.ModelServiceClient.AnalyzeAnomaly(&svc.AnomalyAnalyzeRequest{
			Chart:        req.Chart,
			Metric:       req.Metric,
			SeriesName:   req.SeriesName,
			XName:        req.XName,
			XValue:       req.XValue,
			YName:        req.YName,
			YValue:       req.YValue,
			DataIndex:    req.DataIndex,
			SeriesValues: req.SeriesValues,
			XValues:      req.XValues,
			TimeRange:    req.TimeRange,
			Filters:      req.Filters,
			Extra:        req.Extra,
		})
		if err == nil && modelResp != nil && strings.TrimSpace(modelResp.Analysis) != "" {
			return &types.AnalyzeAnomalyResponse{
				Summary:       modelResp.Analysis,
				Causes:        modelResp.Causes,
				Actions:       modelResp.Actions,
				Confidence:    modelResp.Confidence,
				UsedModel:     true,
				ModelProvider: "ModelService",
			}, nil
		}
	}

	// 2) fallback：简单统计 + 规则解释（保证接口可用）
	return l.ruleBasedFallback(req), nil
}

func (l *AnalyzeAnomalyLogic) ruleBasedFallback(req *types.AnalyzeAnomalyRequest) *types.AnalyzeAnomalyResponse {
	metric := req.Metric
	if strings.TrimSpace(metric) == "" {
		metric = "指标"
	}

	series := req.SeriesValues
	idx := req.DataIndex
	y := req.YValue

	// 默认结论（无序列上下文时）
	resp := &types.AnalyzeAnomalyResponse{
		Summary:      fmt.Sprintf("已收到点击点位：%s=%v。当前未接入模型服务的异常分析端点，返回基于规则的初步排查建议。", metric, y),
		Causes:       []string{"数据口径/统计周期变更", "数据延迟/补数导致尖峰或塌陷", "上游数据源缺失/重复上报"},
		Actions:      []string{"对比前后相邻时间点与同口径指标（如订单量/支付成功数/访客数）", "检查ETL/埋点/采集任务是否在该时间点报错或重跑", "检查是否存在活动、投放、价格、库存等业务变更"},
		Confidence:   0.35,
		UsedModel:    false,
		ModelProvider: "rule-based",
	}

	if len(series) < 5 || idx < 0 || idx >= len(series) {
		return resp
	}

	// 取窗口：优先用“剔除当前点”的全序列来估计分布
	others := make([]float64, 0, len(series)-1)
	for i, v := range series {
		if i == idx {
			continue
		}
		others = append(others, v)
	}

	mean, std := meanStd(others)
	median := medianFloat64(others)
	if std <= 1e-9 {
		std = math.Abs(mean) * 0.01
		if std <= 1e-9 {
			std = 1
		}
	}
	z := (series[idx] - mean) / std

	direction := "波动"
	if series[idx] > median {
		direction = "异常上升"
	} else if series[idx] < median {
		direction = "异常下降"
	}

	absz := math.Abs(z)
	conf := 0.45
	if absz >= 3 {
		conf = 0.75
	} else if absz >= 2 {
		conf = 0.6
	}

	resp.Confidence = conf
	resp.Summary = fmt.Sprintf("%s 在 %s 发生%s：当前点=%.4g，历史均值≈%.4g，中位数≈%.4g，标准差≈%.4g，z≈%.2f。以下为优先级从高到低的排查路径（规则推断）。",
		metric,
		coalesce(req.XValue, fmt.Sprintf("index=%d", idx)),
		direction,
		series[idx], mean, median, std, z,
	)

	if direction == "异常上升" {
		resp.Causes = []string{
			"业务侧：促销/活动/投放带来流量与转化提升（核对活动日历、投放计划、渠道来源）",
			"业务侧：价格下降、优惠券叠加、爆品上新导致集中购买",
			"供给侧：补发货/补库存后订单集中释放（若指标为订单/支付类）",
			"数据侧：延迟补数、重复上报、口径变更导致尖峰（核对ETL重跑与去重逻辑）",
		}
		resp.Actions = []string{
			"对比同时间的访客数/加购/下单/支付成功链路指标，确认是否全链路同步上升",
			"按渠道/地区/品类/店铺做维度拆解，定位异常贡献最大的子维度",
			"检查是否发生ETL补数/重跑、消息队列堆积回放、去重规则变更",
		}
		return resp
	}

	if direction == "异常下降" {
		resp.Causes = []string{
			"供给侧：库存缺货/商品下架/价格异常导致转化下降（核对库存与商品状态）",
			"业务侧：投放停止、渠道流量下滑、活动结束导致自然回落",
			"系统侧：下单/支付/库存扣减/搜索等关键链路故障或超时（核对错误率与告警）",
			"数据侧：采集/ETL中断、过滤条件变化、口径变更导致塌陷（核对任务日志与数据源）",
		}
		resp.Actions = []string{
			"对比同时间的错误率/超时/订单创建失败数，排除系统故障",
			"核对库存与商品上架状态，查看缺货率、下架数量是否突增",
			"检查采集/ETL任务是否失败、延迟是否突增、数据源是否缺失",
		}
		return resp
	}

	// 不明显上/下：给更通用的建议
	resp.Causes = []string{
		"该点与历史分布存在偏离，但方向不强：可能是噪声/周期性因素叠加（周末/节假日/结算日）",
		"数据侧轻微异常：延迟补数、抽样变化、去重误差",
		"业务侧结构变化：某些子维度波动互相抵消（需要拆维度确认）",
	}
	resp.Actions = []string{
		"拆解到渠道/品类/店铺/地区，找出贡献最大的子维度",
		"用滑动窗口（前后3-7个点）做局部均值对比，判断是否趋势变化而非离群点",
	}
	return resp
}

func coalesce(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func meanStd(xs []float64) (mean, std float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range xs {
		sum += v
	}
	mean = sum / float64(len(xs))
	var ss float64
	for _, v := range xs {
		d := v - mean
		ss += d * d
	}
	std = math.Sqrt(ss / float64(len(xs)))
	return mean, std
}

func medianFloat64(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	cp := append([]float64(nil), xs...)
	sort.Float64s(cp)
	m := len(cp) / 2
	if len(cp)%2 == 1 {
		return cp[m]
	}
	return (cp[m-1] + cp[m]) / 2
}

