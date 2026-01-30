// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/user"
	"model_mall_backend/backend/internal/svc"
)

// 获取用户地址列表
func GetAddressesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetAddressesLogic(r.Context(), svcCtx)
		resp, err := l.GetAddresses()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
