// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/merchant"
	"model_mall_backend/backend/internal/svc"
)

// 获取商户数据统计
func GetMerchantDashboardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewGetMerchantDashboardLogic(r.Context(), svcCtx)
		resp, err := l.GetMerchantDashboard()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
