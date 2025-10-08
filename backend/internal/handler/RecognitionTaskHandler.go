package handler

import (
	"net/http"

	"model_mall_backend/backend/internal/logic"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetTaskStatusHandler 获取任务状态
func GetTaskStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTaskStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRecognitionTaskLogic(r.Context(), svcCtx)
		resp, err := l.GetTaskStatus(&req)
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

// GetTaskListHandler 获取任务列表
func GetTaskListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTaskListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRecognitionTaskLogic(r.Context(), svcCtx)
		resp, err := l.GetTaskList(&req)
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
