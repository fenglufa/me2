package logic

import (
	"context"
	"fmt"

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
	// 1. 参数验证
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}

	// 设置默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 2. 查询记账列表
	expenses, err := l.svcCtx.NoteExpenseModel.FindByUserId(in.UserId, in.StartDate, in.EndDate, in.Category, page, pageSize)
	if err != nil {
		l.Errorf("查询记账列表失败: %v", err)
		return nil, fmt.Errorf("查询记账列表失败")
	}

	// 3. 查询总数
	total, err := l.svcCtx.NoteExpenseModel.Count(in.UserId, in.StartDate, in.EndDate, in.Category)
	if err != nil {
		l.Errorf("统计记账总数失败: %v", err)
		return nil, fmt.Errorf("统计记账总数失败")
	}

	// 4. 计算总金额
	totalAmount, err := l.svcCtx.NoteExpenseModel.SumAmount(in.UserId, in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("统计总金额失败: %v", err)
		// 不影响主流程
		totalAmount = 0
	}

	// 5. 转换为响应格式
	expenseInfos := make([]*note.ExpenseInfo, 0, len(expenses))
	for _, e := range expenses {
		expenseInfo := &note.ExpenseInfo{
			Id:        e.Id,
			NoteId:    e.NoteId,
			UserId:    e.UserId,
			Item:      e.Item,
			Amount:    e.Amount,
			Category:  e.Category,
			CreatedAt: e.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		expenseInfos = append(expenseInfos, expenseInfo)
	}

	l.Infof("成功查询记账列表 (user_id=%d, page=%d, page_size=%d, total=%d, total_amount=%.2f)",
		in.UserId, page, pageSize, total, totalAmount)

	return &note.GetExpensesResponse{
		Expenses:    expenseInfos,
		Total:       total,
		TotalAmount: totalAmount,
	}, nil
}
