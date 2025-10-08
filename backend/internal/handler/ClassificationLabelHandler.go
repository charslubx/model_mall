package handler

import (
	"net/http"

	"model_mall_backend/backend/internal/logic"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetImageLabelsHandler 获取图片标签
func GetImageLabelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetImageLabelsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewClassificationLabelLogic(r.Context(), svcCtx)
		resp, err := l.GetImageLabels(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    0,
				Message: "success",
				Data:    resp,
			})
		}
	}
}
