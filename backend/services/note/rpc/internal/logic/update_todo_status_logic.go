package logic

import (
	"context"
	"fmt"

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
	// 1. 参数验证
	if in.TodoId == 0 {
		return nil, fmt.Errorf("todo_id 不能为空")
	}
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}
	if in.Status < 0 || in.Status > 1 {
		return nil, fmt.Errorf("status 参数无效")
	}

	// 2. 更新TODO状态
	err := l.svcCtx.NoteTodoModel.UpdateStatus(in.TodoId, in.UserId, in.Status)
	if err != nil {
		l.Errorf("更新TODO状态失败: %v", err)
		return nil, fmt.Errorf("更新TODO状态失败")
	}

	l.Infof("成功更新TODO状态 (todo_id=%d, user_id=%d, status=%d)", in.TodoId, in.UserId, in.Status)

	return &note.UpdateTodoStatusResponse{
		Success: true,
	}, nil
}
