package user

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/user"
	"github.com/me2/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取订阅信息
func GetSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetSubscriptionLogic(r.Context(), svcCtx)
		resp, err := l.GetSubscription()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
