package logic

import (
	"context"
	"fmt"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotesLogic {
	return &GetNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取笔记列表
func (l *GetNotesLogic) GetNotes(in *note.GetNotesRequest) (*note.GetNotesResponse, error) {
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
		pageSize = 100 // 限制最大分页数
	}

	// 2. 查询笔记列表
	notes, err := l.svcCtx.NoteModel.FindByUserId(in.UserId, page, pageSize, in.Types)
	if err != nil {
		l.Errorf("查询笔记列表失败: %v", err)
		return nil, fmt.Errorf("查询笔记列表失败")
	}

	// 3. 查询总数
	total, err := l.svcCtx.NoteModel.CountByUserId(in.UserId, in.Types)
	if err != nil {
		l.Errorf("统计笔记总数失败: %v", err)
		return nil, fmt.Errorf("统计笔记总数失败")
	}

	// 4. 转换为响应格式
	noteInfos := make([]*note.NoteInfo, 0, len(notes))
	for _, n := range notes {
		noteInfo := &note.NoteInfo{
			Id:        n.Id,
			UserId:    n.UserId,
			RawText:   n.RawText,
			AiSummary: n.AiSummary,
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: n.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// 处理可选字段 avatar_id
		if n.AvatarId.Valid {
			noteInfo.AvatarId = n.AvatarId.Int64
		}

		// 解析 types JSON
		if n.Types.Valid && n.Types.String != "" {
			noteInfo.Types = parseTypesJson(n.Types.String)
		}

		// 解析 emotion_data JSON
		if n.EmotionData.Valid && n.EmotionData.String != "" {
			noteInfo.EmotionData = parseEmotionDataJson(n.EmotionData.String)
		}

		noteInfos = append(noteInfos, noteInfo)
	}

	l.Infof("成功查询笔记列表 (user_id=%d, page=%d, page_size=%d, total=%d)",
		in.UserId, page, pageSize, total)

	return &note.GetNotesResponse{
		Notes: noteInfos,
		Total: total,
	}, nil
}
