package handler

import (
	"net/http"

	"model_mall_backend/backend/internal/logic"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UploadImageHandler 上传图片
func UploadImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadImageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.UploadImage(r, &req)
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

// GetImageListHandler 获取图片列表
func GetImageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetImageListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.GetImageList(&req)
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
