package diary

import (
	"net/http"

	"github.com/me2/gateway/internal/logic/diary"
	"github.com/me2/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取日记统计
func GetDiaryStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := diary.NewGetDiaryStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetDiaryStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
