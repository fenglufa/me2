package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/me2/pkg/errcode"
	"github.com/me2/pkg/response"
	"github.com/zeromicro/go-zero/core/logx"
)

// Recovery panic 恢复中间件
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// 记录 panic 日志
					logx.Errorw("Panic recovered",
						logx.Field("error", fmt.Sprintf("%v", err)),
						logx.Field("stack", string(debug.Stack())),
						logx.Field("path", r.URL.Path),
						logx.Field("method", r.Method),
					)

					// 返回错误响应
					response.HttpErrorWithStatus(w, http.StatusInternalServerError, errcode.ErrInternalServer)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
