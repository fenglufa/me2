package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTodoStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTodoStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTodoStatusLogic {
	return &UpdateTodoStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新 TODO 状态
func (l *UpdateTodoStatusLogic) UpdateTodoStatus(in *note.UpdateTodoStatusRequest) (*note.UpdateTodoStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &note.UpdateTodoStatusResponse{}, nil
}
