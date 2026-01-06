// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package seller

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/seller"
	"model_mall_backend/backend/internal/svc"
)

// 获取卖家信息
func GetSellerInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := seller.NewGetSellerInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetSellerInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
