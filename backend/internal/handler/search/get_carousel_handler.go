// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package search

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/search"
	"model_mall_backend/backend/internal/svc"
)

// 获取轮播图数据
func GetCarouselHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := search.NewGetCarouselLogic(r.Context(), svcCtx)
		resp, err := l.GetCarousel()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
