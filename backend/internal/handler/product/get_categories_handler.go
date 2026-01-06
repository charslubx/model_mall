// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package product

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/product"
	"model_mall_backend/backend/internal/svc"
)

// 获取商品分类列表
func GetCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.GetCategories()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
