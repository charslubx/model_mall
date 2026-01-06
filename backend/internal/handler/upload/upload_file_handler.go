// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package upload

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/upload"
	"model_mall_backend/backend/internal/svc"
)

// 上传文件
func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := upload.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
