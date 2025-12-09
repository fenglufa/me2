package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取笔记详情
func NewGetNoteDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteDetailLogic {
	return &GetNoteDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoteDetailLogic) GetNoteDetail(req *types.GetNoteDetailRequest) (resp *types.GetNoteDetailResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.GetNoteDetail(l.ctx, &note.GetNoteDetailRequest{
		UserId: userID,
		NoteId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// 转换笔记数据
	noteResp := types.NoteResponse{
		Id:        rpcResp.Note.Id,
		UserId:    rpcResp.Note.UserId,
		AvatarId:  rpcResp.Note.AvatarId,
		RawText:   rpcResp.Note.RawText,
		AiSummary: rpcResp.Note.AiSummary,
		Types:     rpcResp.Note.Types,
		EmotionData: types.EmotionData{
			Primary: rpcResp.Note.EmotionData.Primary,
			Score:   rpcResp.Note.EmotionData.Score,
		},
		CreatedAt: rpcResp.Note.CreatedAt,
		UpdatedAt: rpcResp.Note.UpdatedAt,
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

	// 转换记账列表
	expenses := make([]types.ExpenseResponse, 0, len(rpcResp.Expenses))
	for _, e := range rpcResp.Expenses {
		expenses = append(expenses, types.ExpenseResponse{
			Id:        e.Id,
			NoteId:    e.NoteId,
			UserId:    e.UserId,
			Item:      e.Item,
			Amount:    e.Amount,
			Category:  e.Category,
			CreatedAt: e.CreatedAt,
		})
	}

	return &types.GetNoteDetailResponse{
		Note:     noteResp,
		Todos:    todos,
		Expenses: expenses,
	}, nil
}
