package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"model_mall_backend/backend/internal/logic"
	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/repository"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ClassifyImageHandler 图片分类处理器
func ClassifyImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析multipart表单
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "解析表单失败: " + err.Error(),
			})
			return
		}
		
		// 获取图片文件
		file, header, err := r.FormFile("image")
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "获取图片文件失败: " + err.Error(),
			})
			return
		}
		defer file.Close()
		
		// 读取图片数据
		imageData, err := io.ReadAll(file)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "读取图片数据失败: " + err.Error(),
			})
			return
		}
		
		// 获取其他参数
		modelName := r.FormValue("model_name")
		if modelName == "" {
			modelName = "default"
		}
		
		saveImage := r.FormValue("save_image") == "true"
		
		var minConfidence float64
		if minConfidenceStr := r.FormValue("min_confidence"); minConfidenceStr != "" {
			minConfidence, _ = strconv.ParseFloat(minConfidenceStr, 64)
		}
		
		// 获取用户ID（从JWT中间件获取）
		var userID int64
		if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
			if uid, ok := userIDValue.(int64); ok {
				userID = uid
			}
		}
		
		// 构建请求
		req := &models.ImageClassificationReq{
			ImageData:     imageData,
			ImageName:     header.Filename,
			ModelName:     modelName,
			UserID:        userID,
			SaveImage:     saveImage,
			MinConfidence: minConfidence,
		}
		
		// 调用业务逻辑
		l := logic.NewImageClassificationLogic(r.Context(), svcCtx)
		resp, err := l.ClassifyImage(req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		
		// 返回成功响应
		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    200,
			Message: "分类成功",
			Data:    resp,
		})
	}
}

// GetClassificationHandler 获取分类记录处理器
func GetClassificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取分类ID
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "缺少分类ID参数",
			})
			return
		}
		
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "无效的分类ID",
			})
			return
		}
		
		// 调用业务逻辑
		l := logic.NewImageClassificationLogic(r.Context(), svcCtx)
		resp, err := l.GetClassification(id)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		
		// 返回成功响应
		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    200,
			Message: "获取成功",
			Data:    resp,
		})
	}
}

// GetUserClassificationsHandler 获取用户分类记录列表处理器
func GetUserClassificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户ID（从JWT中间件获取）
		var userID int64
		if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
			if uid, ok := userIDValue.(int64); ok {
				userID = uid
			}
		}
		
		if userID == 0 {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    401,
				Message: "未授权访问",
			})
			return
		}
		
		// 获取分页参数
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page <= 0 {
			page = 1
		}
		
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		if pageSize <= 0 || pageSize > 100 {
			pageSize = 20
		}
		
		// 调用业务逻辑
		l := logic.NewImageClassificationLogic(r.Context(), svcCtx)
		resp, err := l.GetUserClassifications(userID, page, pageSize)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		
		// 返回成功响应
		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    200,
			Message: "获取成功",
			Data:    resp,
		})
	}
}

// GetStatisticsHandler 获取统计信息处理器
func GetStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取用户ID（从JWT中间件获取）
		var userID *int64
		if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
			if uid, ok := userIDValue.(int64); ok && uid > 0 {
				userID = &uid
			}
		}
		
		// 检查是否是管理员权限（可以查看全局统计）
		isAdmin := false
		if roleValue := r.Context().Value("role_code"); roleValue != nil {
			if roleCode, ok := roleValue.(string); ok {
				isAdmin = roleCode == "admin" || roleCode == "super_admin"
			}
		}
		
		// 如果不是管理员，只能查看自己的统计
		if !isAdmin && userID == nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    401,
				Message: "未授权访问",
			})
			return
		}
		
		// 如果是管理员且没有指定用户ID，查看全局统计
		if isAdmin && r.URL.Query().Get("global") == "true" {
			userID = nil
		}
		
		// 调用业务逻辑
		l := logic.NewImageClassificationLogic(r.Context(), svcCtx)
		resp, err := l.GetStatistics(userID)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		
		// 返回成功响应
		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    200,
			Message: "获取成功",
			Data:    resp,
		})
	}
}

// DeleteClassificationHandler 删除分类记录处理器
func DeleteClassificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取分类ID
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "缺少分类ID参数",
			})
			return
		}
		
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    400,
				Message: "无效的分类ID",
			})
			return
		}
		
		// 调用业务逻辑
		l := logic.NewImageClassificationLogic(r.Context(), svcCtx)
		err = l.DeleteClassification(id)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		
		// 返回成功响应
		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    200,
			Message: "删除成功",
		})
	}
}

// SearchClassificationsHandler 搜索分类记录处理器
func SearchClassificationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取搜索参数
		req := &repository.SearchClassificationReq{
			Page:     1,
			PageSize: 20,
		}
		
		// 解析分页参数
		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
			req.Page = page
		}
		
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil && pageSize > 0 && pageSize <= 100 {
			req.PageSize = pageSize
		}
		
		// 解析搜索条件
		if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
			if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
				req.UserID = &userID
			}
		}
		
		req.ModelName = strings.TrimSpace(r.URL.Query().Get("model_name"))
		req.ImageName = strings.TrimSpace(r.URL.Query().Get("image_name"))
		
		if statusStr := r.URL.Query().Get("status"); statusStr != "" {
			if status, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
				statusInt8 := int8(status)
				req.Status = &statusInt8
			}
		}
		
		if minConfidenceStr := r.URL.Query().Get("min_confidence"); minConfidenceStr != "" {
			if minConfidence, err := strconv.ParseFloat(minConfidenceStr, 64); err == nil {
				req.MinConfidence = &minConfidence
			}
		}
		
		// 获取当前用户ID
		var currentUserID int64
		if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
			if uid, ok := userIDValue.(int64); ok {
				currentUserID = uid
			}
		}
		
		// 检查权限：普通用户只能搜索自己的记录
		isAdmin := false
		if roleValue := r.Context().Value("role_code"); roleValue != nil {
			if roleCode, ok := roleValue.(string); ok {
				isAdmin = roleCode == "admin" || roleCode == "super_admin"
			}
		}
		
		if !isAdmin {
			req.UserID = &currentUserID
		}
		
		// 调用业务逻辑
		l := logic.NewImageClassificationLogic(r.Context(), svcCtx)
		resp, err := l.SearchClassifications(req)
		if err != nil {
			httpx.OkJsonCtx(r.Context(), w, types.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		
		// 返回成功响应
		httpx.OkJsonCtx(r.Context(), w, types.Response{
			Code:    200,
			Message: "搜索成功",
			Data:    resp,
		})
	}
}