package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ResponseMiddleware struct {
}

// 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponseMiddleware() *ResponseMiddleware {
	return &ResponseMiddleware{}
}

func (m *ResponseMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 包装 ResponseWriter 以捕获状态码
		ww := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 处理请求
		next(ww, r)

		// 记录响应状态
		logx.Infof("Response status: %d", ww.statusCode)
	}
}

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// 错误响应
func Error(w http.ResponseWriter, code int, msg string) {
	httpx.WriteJson(w, code, Response{
		Code:    code,
		Message: msg,
	})
}
