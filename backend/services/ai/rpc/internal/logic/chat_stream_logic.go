package logic

import (
	"context"

	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/deepseek"
	"github.com/me2/ai/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 对话生成（流式）
func (l *ChatStreamLogic) ChatStream(in *ai.ChatRequest, stream ai.Ai_ChatStreamServer) error {
	systemPrompt, userPrompt, err := l.svcCtx.PromptRenderer.Render(in.PromptTemplate, in.Variables)
	if err != nil {
		return err
	}

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

	return l.svcCtx.DeepseekClient.ChatStream(l.ctx, req, func(content string, done bool) error {
		return stream.Send(&ai.ChatStreamResponse{
			Content: content,
			Done:    done,
		})
	})
}
