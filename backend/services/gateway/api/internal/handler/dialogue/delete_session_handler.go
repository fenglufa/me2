package dialogue

import (
	"net/http"

	"github.com/me2/gateway/api/internal/logic/dialogue"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除会话
func DeleteSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dialogue.NewDeleteSessionLogic(r.Context(), svcCtx)
		resp, err := l.DeleteSession()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
