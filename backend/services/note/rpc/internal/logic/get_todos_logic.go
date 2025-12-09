package logic

import (
	"context"
	"fmt"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTodosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodosLogic {
	return &GetTodosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取 TODO 列表
func (l *GetTodosLogic) GetTodos(in *note.GetTodosRequest) (*note.GetTodosResponse, error) {
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

	// 2. 查询TODO列表
	todos, err := l.svcCtx.NoteTodoModel.FindByUserId(in.UserId, in.Status, page, pageSize)
	if err != nil {
		l.Errorf("查询TODO列表失败: %v", err)
		return nil, fmt.Errorf("查询TODO列表失败")
	}

	// 3. 查询总数
	total, err := l.svcCtx.NoteTodoModel.Count(in.UserId, in.Status)
	if err != nil {
		l.Errorf("统计TODO总数失败: %v", err)
		return nil, fmt.Errorf("统计TODO总数失败")
	}

	// 4. 转换为响应格式
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

	l.Infof("成功查询TODO列表 (user_id=%d, status=%d, page=%d, page_size=%d, total=%d)",
		in.UserId, in.Status, page, pageSize, total)

	return &note.GetTodosResponse{
		Todos: todoInfos,
		Total: total,
	}, nil
}
