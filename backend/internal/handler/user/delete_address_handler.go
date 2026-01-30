// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/user"
	"model_mall_backend/backend/internal/svc"
)

// 删除收货地址
func DeleteAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewDeleteAddressLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAddress(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
