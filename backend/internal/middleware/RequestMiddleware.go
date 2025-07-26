package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RequestMiddleware struct {
}

func NewRequestMiddleware() *RequestMiddleware {
	return &RequestMiddleware{}
}

func (m *RequestMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. 记录请求开始
		logx.Infof("Request started: %s %s", r.Method, r.URL.Path)

		// 2. 读取并记录请求体
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			// 重新设置请求体，因为ReadAll会清空原有的请求体
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if len(bodyBytes) > 0 {
			logx.Infof("Request body: %s", string(bodyBytes))
		}

		// 3. 继续处理请求
		next(w, r)

		// 4. 记录请求处理时间
		logx.Infof("Request completed: %s %s, took: %v",
			r.Method, r.URL.Path, time.Since(start))
	}
}
