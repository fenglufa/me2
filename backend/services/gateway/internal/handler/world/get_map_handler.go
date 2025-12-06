package world

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/world"
	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取世界地图详情
func GetMapHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMapRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := world.NewGetMapLogic(r.Context(), svcCtx)
		resp, err := l.GetMap(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
