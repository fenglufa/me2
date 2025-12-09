package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodosLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取TODO列表
func NewGetTodosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodosLogic {
	return &GetTodosLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTodosLogic) GetTodos(req *types.GetTodosRequest) (resp *types.GetTodosResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.GetTodos(l.ctx, &note.GetTodosRequest{
		UserId:   userID,
		Status:   int32(req.Status),
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	// 转换TODO列表
	todos := make([]types.TodoResponse, 0, len(rpcResp.Todos))
	for _, t := range rpcResp.Todos {
		todos = append(todos, types.TodoResponse{
			Id:          t.Id,
			NoteId:      t.NoteId,
			UserId:      t.UserId,
			Title:       t.Title,
			DueDate:     t.DueDate,
			Status:      int(t.Status),
			CreatedAt:   t.CreatedAt,
			CompletedAt: t.CompletedAt,
		})
	}

	return &types.GetTodosResponse{
		Total: rpcResp.Total,
		List:  todos,
	}, nil
}
