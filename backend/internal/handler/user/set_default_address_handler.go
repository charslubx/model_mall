// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/user"
	"model_mall_backend/backend/internal/svc"
)

// 设置默认地址
func SetDefaultAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewSetDefaultAddressLogic(r.Context(), svcCtx)
		resp, err := l.SetDefaultAddress(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
