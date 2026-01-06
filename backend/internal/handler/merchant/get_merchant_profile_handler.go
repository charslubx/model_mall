// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package merchant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/merchant"
	"model_mall_backend/backend/internal/svc"
)

// 获取商户信息
func GetMerchantProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := merchant.NewGetMerchantProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetMerchantProfile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
