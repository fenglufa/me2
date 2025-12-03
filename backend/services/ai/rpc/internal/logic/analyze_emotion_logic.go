package logic

import (
	"context"
	"encoding/json"
	"strings"
	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AnalyzeEmotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnalyzeEmotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalyzeEmotionLogic {
	return &AnalyzeEmotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AnalyzeEmotionLogic) AnalyzeEmotion(in *ai.AnalyzeEmotionRequest) (*ai.AnalyzeEmotionResponse, error) {
	chatLogic := NewChatLogic(l.ctx, l.svcCtx)
	chatResp, err := chatLogic.Chat(&ai.ChatRequest{
		PromptTemplate: "emotion_analysis",
		Variables:      map[string]string{"text": in.Text},
		UserId:         in.UserId,
	})
	if err != nil {
		return nil, err
	}

	// 清理响应内容，移除可能的 markdown 代码块标记
	content := l.cleanJSONContent(chatResp.Content)

	var result struct {
		Emotion  string  `json:"emotion"`
		Score    float32 `json:"score"`
		Analysis string  `json:"analysis"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		logx.Errorf("解析情绪分析结果失败: %v, 原始内容: %s", err, chatResp.Content)
		return &ai.AnalyzeEmotionResponse{
			Emotion:  "neutral",
			Score:    0.5,
			Analysis: chatResp.Content,
		}, nil
	}

	return &ai.AnalyzeEmotionResponse{
		Emotion:  result.Emotion,
		Score:    result.Score,
		Analysis: result.Analysis,
	}, nil
}

// cleanJSONContent 清理 JSON 内容，移除 markdown 代码块标记
func (l *AnalyzeEmotionLogic) cleanJSONContent(content string) string {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}
