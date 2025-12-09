package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetExpenseStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取记账统计
func NewGetExpenseStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExpenseStatsLogic {
	return &GetExpenseStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetExpenseStatsLogic) GetExpenseStats(req *types.GetExpenseStatsRequest) (resp *types.GetExpenseStatsResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.GetExpenseStats(l.ctx, &note.GetExpenseStatsRequest{
		UserId:    userID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return nil, err
	}

	// 转换分类统计 (map -> array)
	categoryStats := make([]types.CategoryStat, 0, len(rpcResp.CategoryStats))
	for category, amount := range rpcResp.CategoryStats {
		categoryStats = append(categoryStats, types.CategoryStat{
			Category: category,
			Amount:   amount,
		})
	}

	// 转换每日支出
	dailyExpenses := make([]types.DailyExpense, 0, len(rpcResp.DailyExpenses))
	for _, de := range rpcResp.DailyExpenses {
		dailyExpenses = append(dailyExpenses, types.DailyExpense{
			Date:   de.Date,
			Amount: de.Amount,
		})
	}

	return &types.GetExpenseStatsResponse{
		TotalAmount:   rpcResp.TotalAmount,
		CategoryStats: categoryStats,
		DailyExpenses: dailyExpenses,
	}, nil
}
