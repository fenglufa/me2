package diary

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
