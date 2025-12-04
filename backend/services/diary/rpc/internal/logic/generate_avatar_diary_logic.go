package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/me2/ai/rpc/ai_client"
	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/model"
	"github.com/me2/diary/rpc/internal/svc"
	"github.com/me2/event/rpc/event_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateAvatarDiaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateAvatarDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAvatarDiaryLogic {
	return &GenerateAvatarDiaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成分身日记（由 Scheduler 调用）
func (l *GenerateAvatarDiaryLogic) GenerateAvatarDiary(in *diary.GenerateAvatarDiaryRequest) (*diary.GenerateAvatarDiaryResponse, error) {
	date := in.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// 1. 查询今天的所有事件
	events, err := l.svcCtx.EventRpc.GetEventTimeline(l.ctx, &event_client.GetEventTimelineRequest{
		AvatarId: in.AvatarId,
		Page:     1,
		PageSize: 10,  // 限制最多10个事件，避免内容过长导致AI超时
	})
	if err != nil {
		l.Errorf("查询事件失败: %v", err)
		return nil, fmt.Errorf("查询事件失败")
	}

	if len(events.Events) == 0 {
		l.Infof("分身 %d 在 %s 没有事件，跳过日记生成", in.AvatarId, date)
		return nil, fmt.Errorf("今天没有事件")
	}

	// 2. 构建事件摘要（最多取前5个事件）
	var eventSummary strings.Builder
	maxEvents := 5
	if len(events.Events) < maxEvents {
		maxEvents = len(events.Events)
	}
	for i := 0; i < maxEvents; i++ {
		e := events.Events[i]
		eventSummary.WriteString(fmt.Sprintf("%d. %s\n", i+1, e.EventTitle))
	}

	// 3. 调用 AI 生成日记
	// 使用独立的 context 和配置的超时时间
	timeout := l.svcCtx.Config.AIGenerateTimeout
	if timeout == 0 {
		timeout = 30000 // 默认 30 秒
	}
	aiCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	aiResp, err := l.svcCtx.AIRpc.Generate(aiCtx, &ai_client.GenerateRequest{
		PromptTemplate: "avatar_diary",
		Variables: map[string]string{
			"avatar_name": "我",  // 可以后续从 Avatar 服务获取真实昵称
			"date":        date,
			"events":      eventSummary.String(),
			"emotion":     "平静",  // 可以根据事件类型推断
			"personality": "温暖、善良、乐观",  // 可以从 Avatar 服务获取
		},
		AvatarId: in.AvatarId,
	})
	if err != nil {
		l.Errorf("AI生成日记失败: %v", err)
		return nil, fmt.Errorf("AI生成日记失败")
	}

	title, content, mood, tags := parseDiaryResponse(aiResp.Content)

	// 4. 保存日记
	diaryDate, _ := time.Parse("2006-01-02", date)
	diaryModel := &model.Diaries{
		AvatarId:    in.AvatarId,
		Type:        "avatar",
		Date:        diaryDate,
		Title:       title,
		Content:     content,
		Mood:        mood,
		Tags:        strings.Join(tags, ","),
		IsImportant: 0,
	}

	result, err := l.svcCtx.DiaryModel.Insert(context.Background(), diaryModel)
	if err != nil {
		l.Errorf("保存日记失败: %v", err)
		return nil, fmt.Errorf("保存日记失败")
	}

	diaryId, _ := result.LastInsertId()
	l.Infof("成功生成分身日记 (diary_id=%d, avatar_id=%d, date=%s)", diaryId, in.AvatarId, date)

	return &diary.GenerateAvatarDiaryResponse{
		DiaryId: diaryId,
		Title:   title,
		Content: content,
		Mood:    mood,
		Tags:    tags,
	}, nil
}
