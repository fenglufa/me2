package action

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/action"
	"github.com/me2/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取最近一次行动
func GetLastActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := action.NewGetLastActionLogic(r.Context(), svcCtx)
		resp, err := l.GetLastAction()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
