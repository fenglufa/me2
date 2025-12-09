package logic

import (
	"context"
	"fmt"

	"github.com/me2/note/rpc/internal/model"
	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoteDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteDetailLogic {
	return &GetNoteDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取笔记详情
func (l *GetNoteDetailLogic) GetNoteDetail(in *note.GetNoteDetailRequest) (*note.GetNoteDetailResponse, error) {
	// 1. 参数验证
	if in.NoteId == 0 {
		return nil, fmt.Errorf("note_id 不能为空")
	}
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}

	// 2. 查询笔记
	n, err := l.svcCtx.NoteModel.FindOne(in.NoteId)
	if err != nil {
		l.Errorf("查询笔记失败: %v", err)
		return nil, fmt.Errorf("笔记不存在")
	}

	// 3. 验证笔记所有权
	if n.UserId != in.UserId {
		return nil, fmt.Errorf("无权访问此笔记")
	}

	// 4. 构建笔记信息
	noteInfo := &note.NoteInfo{
		Id:        n.Id,
		UserId:    n.UserId,
		RawText:   n.RawText,
		AiSummary: n.AiSummary,
		CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: n.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if n.AvatarId.Valid {
		noteInfo.AvatarId = n.AvatarId.Int64
	}

	if n.Types.Valid && n.Types.String != "" {
		noteInfo.Types = parseTypesJson(n.Types.String)
	}

	if n.EmotionData.Valid && n.EmotionData.String != "" {
		noteInfo.EmotionData = parseEmotionDataJson(n.EmotionData.String)
	}

	// 5. 查询关联的 TODO
	todos, err := l.svcCtx.NoteTodoModel.FindByNoteId(in.NoteId)
	if err != nil {
		l.Errorf("查询TODO失败: %v", err)
		// 不影响主流程，返回空列表
		todos = []*model.NoteTodo{}
	}

	todoInfos := make([]*note.TodoInfo, 0, len(todos))
	for _, t := range todos {
		todoInfo := &note.TodoInfo{
			Id:        t.Id,
			NoteId:    t.NoteId,
			UserId:    t.UserId,
			Title:     t.Title,
			Status:    t.Status,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if t.DueDate.Valid {
			todoInfo.DueDate = t.DueDate.String
		}

		if t.CompletedAt.Valid {
			todoInfo.CompletedAt = t.CompletedAt.Time.Format("2006-01-02 15:04:05")
		}

		todoInfos = append(todoInfos, todoInfo)
	}

	// 6. 查询关联的记账记录
	expenses, err := l.svcCtx.NoteExpenseModel.FindByNoteId(in.NoteId)
	if err != nil {
		l.Errorf("查询记账记录失败: %v", err)
		// 不影响主流程，返回空列表
		expenses = []*model.NoteExpense{}
	}

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

	l.Infof("成功查询笔记详情 (note_id=%d, user_id=%d, todos=%d, expenses=%d)",
		in.NoteId, in.UserId, len(todoInfos), len(expenseInfos))

	return &note.GetNoteDetailResponse{
		Note:     noteInfo,
		Todos:    todoInfos,
		Expenses: expenseInfos,
	}, nil
}
