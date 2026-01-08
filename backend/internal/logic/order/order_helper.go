package order

import "model_mall_backend/backend/internal/models"

// OrderStatusToString 将订单状态整数转换为字符串
func OrderStatusToString(status int) string {
	statusMap := map[int]string{
		models.OrderStatusPending:   "pending",
		models.OrderStatusPaid:      "paid",
		models.OrderStatusShipped:   "shipped",
		models.OrderStatusCompleted: "completed",
		models.OrderStatusCancelled: "cancelled",
	}
	if s, ok := statusMap[status]; ok {
		return s
	}
	return "unknown"
}

// OrderStatusTextMap 订单状态文本映射
var OrderStatusTextMap = map[string]string{
	"pending":   "待付款",
	"paid":      "已付款",
	"shipped":   "已发货",
	"completed": "已完成",
	"cancelled": "已取消",
}

// PaymentMethodTextMap 支付方式文本映射
var PaymentMethodTextMap = map[string]string{
	"alipay": "支付宝",
	"wechat": "微信支付",
	"union":  "银联支付",
}
