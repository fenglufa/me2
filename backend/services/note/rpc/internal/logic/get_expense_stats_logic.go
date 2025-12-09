package logic

import (
	"context"
	"fmt"

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
	// 1. 参数验证
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}

	// 2. 计算总金额
	totalAmount, err := l.svcCtx.NoteExpenseModel.SumAmount(in.UserId, in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("统计总金额失败: %v", err)
		return nil, fmt.Errorf("统计总金额失败")
	}

	// 3. 按分类统计
	categoryStats, err := l.svcCtx.NoteExpenseModel.StatsByCategory(in.UserId, in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("按分类统计失败: %v", err)
		return nil, fmt.Errorf("按分类统计失败")
	}

	// 4. 计算每日支出（简单实现，返回空列表）
	// TODO: 实现每日支出统计
	dailyExpenses := make([]*note.DailyExpense, 0)

	l.Infof("成功查询记账统计 (user_id=%d, total_amount=%.2f, categories=%d)",
		in.UserId, totalAmount, len(categoryStats))

	return &note.GetExpenseStatsResponse{
		TotalAmount:    totalAmount,
		CategoryStats:  categoryStats,
		DailyExpenses:  dailyExpenses,
	}, nil
}
