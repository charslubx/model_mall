// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/merchant"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"
)

// 获取商户数据分析
func GetMerchantAnalyticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAnalyticsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := merchant.NewGetMerchantAnalyticsLogic(r.Context(), svcCtx)
		resp, err := l.GetMerchantAnalytics(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
