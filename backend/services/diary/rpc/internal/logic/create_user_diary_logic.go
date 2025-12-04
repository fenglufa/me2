package logic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/me2/ai/rpc/ai_client"
	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/model"
	"github.com/me2/diary/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserDiaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserDiaryLogic {
	return &CreateUserDiaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建用户日记
func (l *CreateUserDiaryLogic) CreateUserDiary(in *diary.CreateUserDiaryRequest) (*diary.CreateUserDiaryResponse, error) {
	// 1. 检查今天是否已有用户日记
	today := time.Now().Format("2006-01-02")
	existingDiaries, err := l.svcCtx.DiaryModel.FindByAvatarAndType(l.ctx, in.AvatarId, "user", 1, 1, today, today)
	if err != nil {
		l.Errorf("查询现有日记失败: %v", err)
	}

	// 如果今天已经有日记，返回错误提示
	if len(existingDiaries) > 0 {
		return nil, fmt.Errorf("今天已经写过日记了，一天只能写一篇用户日记")
	}

	// 2. 保存用户日记
	isImportant := int64(0)
	if in.IsImportant {
		isImportant = 1
	}

	diaryModel := &model.Diaries{
		AvatarId:    in.AvatarId,
		Type:        "user",
		Date:        time.Now(),
		Title:       in.Title,
		Content:     in.Content,
		Tags:        strings.Join(in.Tags, ","),
		IsImportant: isImportant,
	}

	result, err := l.svcCtx.DiaryModel.Insert(context.Background(), diaryModel)
	if err != nil {
		l.Errorf("保存用户日记失败: %v", err)
		return nil, fmt.Errorf("保存日记失败")
	}

	diaryId, _ := result.LastInsertId()

	// 2. 调用 AI 进行情绪分析
	// 使用独立的 context 和配置的超时时间
	timeout := l.svcCtx.Config.AIGenerateTimeout
	if timeout == 0 {
		timeout = 30000 // 默认 30 秒
	}
	aiCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	emotionResp, err := l.svcCtx.AIRpc.Generate(aiCtx, &ai_client.GenerateRequest{
		PromptTemplate: "emotion_analysis",
		Variables: map[string]string{
			"text": in.Content,
		},
		AvatarId: in.AvatarId,
	})
	if err != nil {
		l.Errorf("情绪分析失败: %v", err)
	}

	emotionScore := int32(0)
	if emotionResp != nil {
		emotionScore = parseEmotionScore(emotionResp.Content)
	}

	// 3. 调用 AI 生成分身回应 (使用同一个 context)
	replyResp, err := l.svcCtx.AIRpc.Generate(aiCtx, &ai_client.GenerateRequest{
		PromptTemplate: "user_diary_reply",
		Variables: map[string]string{
			"diary_content": in.Content,
			"emotion":       fmt.Sprintf("情绪分数: %d", emotionScore),
			"personality":   "温暖、善良、乐观",  // 可以从 Avatar 服务获取
		},
		AvatarId: in.AvatarId,
	})
	if err != nil {
		l.Errorf("生成分身回应失败: %v", err)
	}

	replyContent := ""
	if replyResp != nil {
		replyContent = replyResp.Content
	}

	// 4. 更新日记的回应和情绪分数
	err = l.svcCtx.DiaryModel.Update(context.Background(), &model.Diaries{
		Id:           diaryId,
		AvatarId:     in.AvatarId,
		Type:         "user",
		Date:         time.Now(),
		Title:        in.Title,
		Content:      in.Content,
		Tags:         strings.Join(in.Tags, ","),
		ReplyContent: sql.NullString{String: replyContent, Valid: true},
		EmotionScore: int64(emotionScore),
		IsImportant:  isImportant,
	})
	if err != nil {
		l.Errorf("更新日记失败: %v", err)
	}

	l.Infof("成功创建用户日记 (diary_id=%d, avatar_id=%d)", diaryId, in.AvatarId)

	return &diary.CreateUserDiaryResponse{
		DiaryId:      diaryId,
		ReplyContent: replyContent,
		EmotionScore: emotionScore,
	}, nil
}
