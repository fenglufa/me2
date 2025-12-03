package logic

import (
	"context"
	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type EmbeddingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmbeddingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmbeddingLogic {
	return &EmbeddingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmbeddingLogic) Embedding(in *ai.EmbeddingRequest) (*ai.EmbeddingResponse, error) {
	// TODO: 实现向量化功能
	return &ai.EmbeddingResponse{
		Vector:    []float32{},
		Dimension: 0,
	}, nil
}
