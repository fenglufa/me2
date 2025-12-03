package logic

import (
	"context"

	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromptsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPromptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromptsLogic {
	return &GetPromptsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取 Prompt 模板列表
func (l *GetPromptsLogic) GetPrompts(in *ai.GetPromptsRequest) (*ai.GetPromptsResponse, error) {
	// todo: add your logic here and delete this line

	return &ai.GetPromptsResponse{}, nil
}
