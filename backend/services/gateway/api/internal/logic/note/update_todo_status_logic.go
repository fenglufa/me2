package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTodoStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新TODO状态
func NewUpdateTodoStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTodoStatusLogic {
	return &UpdateTodoStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTodoStatusLogic) UpdateTodoStatus(req *types.UpdateTodoStatusRequest) (resp *types.UpdateTodoStatusResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.UpdateTodoStatus(l.ctx, &note.UpdateTodoStatusRequest{
		UserId: userID,
		TodoId: req.Id,
		Status: int32(req.Status),
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateTodoStatusResponse{
		Success: rpcResp.Success,
	}, nil
}
