package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoteLogic {
	return &UpdateNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新笔记
func (l *UpdateNoteLogic) UpdateNote(in *note.UpdateNoteRequest) (*note.UpdateNoteResponse, error) {
	// 1. 参数验证
	if in.NoteId == 0 {
		return nil, fmt.Errorf("note_id 不能为空")
	}
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}
	if in.RawText == "" {
		return nil, fmt.Errorf("笔记内容不能为空")
	}

	// 2. 查询笔记验证所有权
	existingNote, err := l.svcCtx.NoteModel.FindOne(in.NoteId)
	if err != nil {
		l.Errorf("查询笔记失败: %v", err)
		return nil, fmt.Errorf("笔记不存在")
	}

	if existingNote.UserId != in.UserId {
		return nil, fmt.Errorf("无权修改此笔记")
	}

	// 3. MVP 阶段：使用简单分类（TODO: Phase 2 使用 AI）
	types, emotionData := simpleClassify(in.RawText)

	// 4. 构建更新对象
	existingNote.RawText = in.RawText
	existingNote.Types = sql.NullString{String: types, Valid: true}
	existingNote.EmotionData = sql.NullString{String: emotionData, Valid: true}

	// 5. 更新笔记
	err = l.svcCtx.NoteModel.Update(existingNote)
	if err != nil {
		l.Errorf("更新笔记失败: %v", err)
		return nil, fmt.Errorf("更新笔记失败")
	}

	l.Infof("成功更新笔记 (note_id=%d, user_id=%d)", in.NoteId, in.UserId)

	return &note.UpdateNoteResponse{
		Success:   true,
		AiSummary: "", // MVP 阶段暂不生成摘要
		Types:     parseTypesJson(types),
	}, nil
}

// simpleClassify MVP 阶段的简单分类（基于关键词）
func simpleClassify(text string) (types string, emotionData string) {
	types = `["note"]`
	emotionData = `{"primary":"neutral","score":0.5}`
	return types, emotionData
}

