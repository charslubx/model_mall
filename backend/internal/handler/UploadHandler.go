package handler

import (
	"net/http"

	"model_mall_backend/backend/internal/logic"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UploadFileHandler 上传通用文件（<= Config.Upload.MaxSize）
func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxSize := svcCtx.Config.Upload.MaxSize
		if maxSize <= 0 {
			maxSize = 10 * 1024 * 1024
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)

		l := logic.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	}
}

// UploadPropConfigHandler 上传属性配置文件（<=10MB）
func UploadPropConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const maxSize = int64(10 * 1024 * 1024)
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)

		l := logic.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.UploadPropConfig(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	}
}
