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

// 更新商户信息
func UpdateMerchantProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMerchantProfileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := merchant.NewUpdateMerchantProfileLogic(r.Context(), svcCtx)
		resp, err := l.UpdateMerchantProfile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
