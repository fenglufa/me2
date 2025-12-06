package avatar

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/avatar"
	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取分身详情
func GetAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAvatarRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := avatar.NewGetAvatarLogic(r.Context(), svcCtx)
		resp, err := l.GetAvatar(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
