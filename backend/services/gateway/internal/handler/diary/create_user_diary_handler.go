package diary

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/diary"
	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建用户日记
func CreateUserDiaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateUserDiaryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := diary.NewCreateUserDiaryLogic(r.Context(), svcCtx)
		resp, err := l.CreateUserDiary(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
