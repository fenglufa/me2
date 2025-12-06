package user

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/user"
	"github.com/me2/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取头像上传凭证
func GetAvatarTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetAvatarTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetAvatarToken()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
