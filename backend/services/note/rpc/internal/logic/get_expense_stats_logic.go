package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExpenseStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExpenseStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExpenseStatsLogic {
	return &GetExpenseStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取记账统计
func (l *GetExpenseStatsLogic) GetExpenseStats(in *note.GetExpenseStatsRequest) (*note.GetExpenseStatsResponse, error) {
	// todo: add your logic here and delete this line

	return &note.GetExpenseStatsResponse{}, nil
}
