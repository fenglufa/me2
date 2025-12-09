package note

import (
	"net/http"

	"github.com/me2/gateway/api/internal/logic/note"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取笔记详情
func GetNoteDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetNoteDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := note.NewGetNoteDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetNoteDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
