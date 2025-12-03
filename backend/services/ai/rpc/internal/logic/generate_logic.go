package logic

import (
	"context"
	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateLogic {
	return &GenerateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateLogic) Generate(in *ai.GenerateRequest) (*ai.GenerateResponse, error) {
	chatLogic := NewChatLogic(l.ctx, l.svcCtx)
	chatResp, err := chatLogic.Chat(&ai.ChatRequest{
		PromptTemplate: in.PromptTemplate,
		Variables:      in.Variables,
		ModelConfig:    in.ModelConfig,
		UserId:         in.UserId,
		AvatarId:       in.AvatarId,
	})
	if err != nil {
		return nil, err
	}

	return &ai.GenerateResponse{
		Content:      chatResp.Content,
		InputTokens:  chatResp.InputTokens,
		OutputTokens: chatResp.OutputTokens,
		DurationMs:   chatResp.DurationMs,
	}, nil
}
