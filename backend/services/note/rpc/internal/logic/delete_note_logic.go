package logic

import (
	"context"
	"fmt"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoteLogic {
	return &DeleteNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除笔记
func (l *DeleteNoteLogic) DeleteNote(in *note.DeleteNoteRequest) (*note.DeleteNoteResponse, error) {
	// 1. 参数验证
	if in.NoteId == 0 {
		return nil, fmt.Errorf("note_id 不能为空")
	}
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}

	// 2. 查询笔记验证所有权
	existingNote, err := l.svcCtx.NoteModel.FindOne(in.NoteId)
	if err != nil {
		l.Errorf("查询笔记失败: %v", err)
		return nil, fmt.Errorf("笔记不存在")
	}

	if existingNote.UserId != in.UserId {
		return nil, fmt.Errorf("无权删除此笔记")
	}

	// 3. 删除关联的 TODO
	if err := l.svcCtx.NoteTodoModel.DeleteByNoteId(in.NoteId); err != nil {
		l.Errorf("删除关联TODO失败: %v", err)
		// 继续执行，不中断
	}

	// 4. 删除关联的记账记录
	if err := l.svcCtx.NoteExpenseModel.DeleteByNoteId(in.NoteId); err != nil {
		l.Errorf("删除关联记账记录失败: %v", err)
		// 继续执行，不中断
	}

	// 5. 删除笔记
	if err := l.svcCtx.NoteModel.Delete(in.NoteId, in.UserId); err != nil {
		l.Errorf("删除笔记失败: %v", err)
		return nil, fmt.Errorf("删除笔记失败")
	}

	l.Infof("成功删除笔记 (note_id=%d, user_id=%d)", in.NoteId, in.UserId)

	return &note.DeleteNoteResponse{
		Success: true,
	}, nil
}

