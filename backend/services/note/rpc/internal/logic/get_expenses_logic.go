package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExpensesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExpensesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExpensesLogic {
	return &GetExpensesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取记账列表
func (l *GetExpensesLogic) GetExpenses(in *note.GetExpensesRequest) (*note.GetExpensesResponse, error) {
	// todo: add your logic here and delete this line

	return &note.GetExpensesResponse{}, nil
}
