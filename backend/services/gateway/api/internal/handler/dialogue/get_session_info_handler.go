package dialogue

import (
	"net/http"

	"github.com/me2/gateway/api/internal/logic/dialogue"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取会话详情
func GetSessionInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSessionInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dialogue.NewGetSessionInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetSessionInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
