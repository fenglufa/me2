package logic

import (
	"context"

	"github.com/me2/ai/rpc/ai"
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
	// todo: add your logic here and delete this line

	return nil
}
