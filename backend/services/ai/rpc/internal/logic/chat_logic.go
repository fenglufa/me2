package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/deepseek"
	"github.com/me2/ai/rpc/internal/model"
	"github.com/me2/ai/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatLogic) Chat(in *ai.ChatRequest) (*ai.ChatResponse, error) {
	startTime := time.Now()

	// 渲染 Prompt
	systemPrompt, userPrompt, err := l.svcCtx.PromptRenderer.Render(in.PromptTemplate, in.Variables)
	if err != nil {
		l.logError(in, 0, err, time.Since(startTime))
		return nil, fmt.Errorf("渲染Prompt失败: %w", err)
	}

	// 构建请求
	req := &deepseek.ChatRequest{
		Model: l.svcCtx.Config.Deepseek.Model,
		Messages: []deepseek.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	if in.ModelConfig != nil {
		req.Temperature = in.ModelConfig.Temperature
		req.MaxTokens = in.ModelConfig.MaxTokens
		req.TopP = in.ModelConfig.TopP
	}

	// 调用 Deepseek
	resp, err := l.svcCtx.DeepseekClient.Chat(l.ctx, req)
	if err != nil {
		l.logError(in, 0, err, time.Since(startTime))
		return nil, fmt.Errorf("调用Deepseek失败: %w", err)
	}

	duration := time.Since(startTime)

	// 记录日志
	l.logSuccess(in, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, duration)

	return &ai.ChatResponse{
		Content:      resp.Choices[0].Message.Content,
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		DurationMs:   duration.Milliseconds(),
	}, nil
}

func (l *ChatLogic) logSuccess(in *ai.ChatRequest, inputTokens, outputTokens int32, duration time.Duration) {
	cost := calculateCost(inputTokens, outputTokens)
	log := &model.AiCallLog{
		ServiceName:  "ai",
		SceneType:    in.PromptTemplate,
		UserId:       in.UserId,
		AvatarId:     in.AvatarId,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		Cost:         cost,
		DurationMs:   duration.Milliseconds(),
		Status:       "success",
	}
	l.svcCtx.LogModel.Insert(log)
}

func (l *ChatLogic) logError(in *ai.ChatRequest, tokens int32, err error, duration time.Duration) {
	log := &model.AiCallLog{
		ServiceName:  "ai",
		SceneType:    in.PromptTemplate,
		UserId:       in.UserId,
		AvatarId:     in.AvatarId,
		InputTokens:  tokens,
		DurationMs:   duration.Milliseconds(),
		Status:       "error",
		ErrorMessage: err.Error(),
	}
	l.svcCtx.LogModel.Insert(log)
}

func calculateCost(inputTokens, outputTokens int32) int64 {
	// Deepseek 价格: 输入 0.001元/1K tokens, 输出 0.002元/1K tokens
	// 转换为分
	inputCost := float64(inputTokens) / 1000.0 * 0.1   // 0.001元 = 0.1分
	outputCost := float64(outputTokens) / 1000.0 * 0.2 // 0.002元 = 0.2分
	return int64(inputCost + outputCost)
}
