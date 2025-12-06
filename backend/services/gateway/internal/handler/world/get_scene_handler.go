package world

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/world"
	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取场景详情
func GetSceneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSceneRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := world.NewGetSceneLogic(r.Context(), svcCtx)
		resp, err := l.GetScene(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
