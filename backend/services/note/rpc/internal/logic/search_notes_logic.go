package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchNotesLogic {
	return &SearchNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索笔记（用于对话集成）
func (l *SearchNotesLogic) SearchNotes(in *note.SearchNotesRequest) (*note.SearchNotesResponse, error) {
	// 1. 参数验证
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}

	// 设置默认limit
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 2. 解析时间范围
	startTime, endTime := parseDateRange(in.DateRange)

	// 3. 搜索笔记
	notes, err := l.svcCtx.NoteModel.Search(in.UserId, in.Query, in.Types, startTime, endTime, limit)
	if err != nil {
		l.Errorf("搜索笔记失败: %v", err)
		return nil, fmt.Errorf("搜索笔记失败")
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

		if n.AvatarId.Valid {
			noteInfo.AvatarId = n.AvatarId.Int64
		}

		if n.Types.Valid && n.Types.String != "" {
			noteInfo.Types = parseTypesJson(n.Types.String)
		}

		if n.EmotionData.Valid && n.EmotionData.String != "" {
			noteInfo.EmotionData = parseEmotionDataJson(n.EmotionData.String)
		}

		noteInfos = append(noteInfos, noteInfo)
	}

	l.Infof("成功搜索笔记 (user_id=%d, query=%s, count=%d)", in.UserId, in.Query, len(noteInfos))

	return &note.SearchNotesResponse{
		Notes: noteInfos,
	}, nil
}

// parseDateRange 解析日期范围
func parseDateRange(dateRange string) (time.Time, time.Time) {
	now := time.Now()
	var startTime, endTime time.Time

	switch dateRange {
	case "today":
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		startTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
		endTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, now.Location())
	case "last_week":
		startTime = now.AddDate(0, 0, -7)
		endTime = now
	case "last_month":
		startTime = now.AddDate(0, -1, 0)
		endTime = now
	default:
		// 不限制时间范围
		startTime = time.Time{}
		endTime = time.Time{}
	}

	return startTime, endTime
}

