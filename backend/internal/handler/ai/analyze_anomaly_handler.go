package ai

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"model_mall_backend/backend/internal/logic/ai"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"
)

// AnalyzeAnomalyHandler 点击图表点位后进行异常原因分析
func AnalyzeAnomalyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AnalyzeAnomalyRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ai.NewAnalyzeAnomalyLogic(r.Context(), svcCtx)
		resp, err := l.AnalyzeAnomaly(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

