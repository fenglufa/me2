package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTodosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodosLogic {
	return &GetTodosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取 TODO 列表
func (l *GetTodosLogic) GetTodos(in *note.GetTodosRequest) (*note.GetTodosResponse, error) {
	// todo: add your logic here and delete this line

	return &note.GetTodosResponse{}, nil
}
