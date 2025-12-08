package dialogue

import (
	"net/http"

	"github.com/me2/gateway/api/internal/logic/dialogue"
	"github.com/me2/gateway/api/internal/middleware"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StreamHandlerWithAuth(svcCtx *svc.ServiceContext) http.HandlerFunc {
	authMiddleware := middleware.NewAuthMiddleware(svcCtx.Config.Auth.AccessSecret)
	return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		l := dialogue.NewStreamLogic(r.Context(), svcCtx)
		err := l.Stream(w, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	})
}
