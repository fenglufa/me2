package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// responseWriter 包装 ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger 日志中间件
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 包装 ResponseWriter
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// 执行请求
			next.ServeHTTP(rw, r)

			// 记录日志
			duration := time.Since(start)
			logx.Infow("HTTP Request",
				logx.Field("method", r.Method),
				logx.Field("path", r.URL.Path),
				logx.Field("status", rw.statusCode),
				logx.Field("duration_ms", duration.Milliseconds()),
				logx.Field("remote_addr", r.RemoteAddr),
				logx.Field("user_agent", r.UserAgent()),
			)
		})
	}
}
