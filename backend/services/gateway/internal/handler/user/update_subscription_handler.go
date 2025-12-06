package user

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/user"
	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新订阅
func UpdateSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateSubscriptionLogic(r.Context(), svcCtx)
		resp, err := l.UpdateSubscription(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
