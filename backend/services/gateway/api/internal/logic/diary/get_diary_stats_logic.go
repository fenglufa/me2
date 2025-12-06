package diary

import (
	"context"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiaryStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取日记统计
func NewGetDiaryStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiaryStatsLogic {
	return &GetDiaryStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDiaryStatsLogic) GetDiaryStats() (resp *types.DiaryStatsResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DiaryRpc.GetDiaryStats(l.ctx, &diary.GetDiaryStatsRequest{
		AvatarId: userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.DiaryStatsResponse{
		TotalCount:  rpcResp.TotalDiaries,
		AvatarCount: rpcResp.AvatarDiaries,
		UserCount:   rpcResp.UserDiaries,
	}, nil
}
