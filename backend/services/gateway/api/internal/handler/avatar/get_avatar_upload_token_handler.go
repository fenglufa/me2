package avatar

import (
	"net/http"

	"github.com/me2/gateway/api/internal/logic/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取分身头像上传凭证
func GetAvatarUploadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := avatar.NewGetAvatarUploadTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetAvatarUploadToken()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
