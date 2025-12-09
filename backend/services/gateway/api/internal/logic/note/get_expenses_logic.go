package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExpensesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取记账列表
func NewGetExpensesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExpensesLogic {
	return &GetExpensesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExpensesLogic) GetExpenses(req *types.GetExpensesRequest) (resp *types.GetExpensesResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.GetExpenses(l.ctx, &note.GetExpensesRequest{
		UserId:    userID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Category:  req.Category,
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
	})
	if err != nil {
		return nil, err
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

	return &types.GetExpensesResponse{
		Total:       rpcResp.Total,
		TotalAmount: rpcResp.TotalAmount,
		List:        expenses,
	}, nil
}
