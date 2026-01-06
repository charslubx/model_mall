package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
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

		// 2. 读取并记录请求体（避免读取大文件/多段表单，防止占用过多内存）
		var bodyBytes []byte
		if r.Body != nil {
			contentType := r.Header.Get("Content-Type")
			// multipart/form-data 往往包含文件，不能在中间件里读完整 body
			if strings.HasPrefix(contentType, "multipart/form-data") {
				logx.Infof("Request body skipped (multipart/form-data)")
			} else {
				// 仅在 Content-Length 可控且较小的时候读取，避免大包导致内存膨胀
				const maxLogBodySize = 1 << 20 // 1MB
				if r.ContentLength > 0 && r.ContentLength <= maxLogBodySize {
					bodyBytes, _ = io.ReadAll(r.Body)
					// 重新设置请求体，因为ReadAll会清空原有的请求体
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				} else if r.ContentLength == 0 {
					// 无请求体
				} else {
					logx.Infof("Request body skipped (Content-Length=%d)", r.ContentLength)
				}
			}
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
