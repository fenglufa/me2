package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/me2/note/rpc/internal/model"
	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoteLogic {
	return &CreateNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建笔记
func (l *CreateNoteLogic) CreateNote(in *note.CreateNoteRequest) (*note.CreateNoteResponse, error) {
	// 1. 参数验证
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}
	if in.RawText == "" {
		return nil, fmt.Errorf("笔记内容不能为空")
	}

	// 2. 构建笔记对象
	noteModel := &model.Note{
		UserId:  in.UserId,
		RawText: in.RawText,
	}

	if in.AvatarId > 0 {
		noteModel.AvatarId = sql.NullInt64{Int64: in.AvatarId, Valid: true}
	}

	// 3. MVP 阶段：暂时不调用 AI 分类，使用简单的关键词匹配
	// TODO: Phase 2 集成 AI Service 进行智能分类
	types, emotionData := l.simpleClassify(in.RawText)

	noteModel.Types = sql.NullString{String: types, Valid: true}
	noteModel.EmotionData = sql.NullString{String: emotionData, Valid: true}
	noteModel.AiSummary = "" // MVP 阶段暂不生成摘要

	// 4. 插入数据库
	result, err := l.svcCtx.NoteModel.Insert(noteModel)
	if err != nil {
		l.Errorf("插入笔记失败: %v", err)
		return nil, fmt.Errorf("创建笔记失败")
	}

	noteId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取笔记ID失败: %v", err)
		return nil, fmt.Errorf("创建笔记失败")
	}

	l.Infof("成功创建笔记 (note_id=%d, user_id=%d)", noteId, in.UserId)

	// 5. 返回结果
	return &note.CreateNoteResponse{
		NoteId:    noteId,
		AiSummary: noteModel.AiSummary,
		Types:     parseTypes(types),
		EmotionData: &note.EmotionData{
			Primary: "neutral",
			Score:   0.5,
		},
	}, nil
}

// simpleClassify MVP 阶段的简单分类（基于关键词）
// TODO: Phase 2 替换为 AI Service 调用
func (l *CreateNoteLogic) simpleClassify(text string) (types string, emotionData string) {
	// 简单的关键词匹配
	types = `["note"]` // 默认类型
	emotionData = `{"primary":"neutral","score":0.5}`

	// TODO: 实现更复杂的关键词匹配逻辑
	// 例如：检测 "花了", "元" -> expense
	//       检测 "要做", "任务" -> todo
	//       检测 "开心", "难过" -> emotion

	return types, emotionData
}

// parseTypes 解析类型字符串为数组
func parseTypes(typesJson string) []string {
	// 简单实现，TODO: 使用 JSON 解析
	if typesJson == "" {
		return []string{}
	}
	return []string{"note"}
}
