package logic

import (
	"context"

	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHealthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthCheckLogic {
	return &HealthCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 健康检查
func (l *HealthCheckLogic) HealthCheck(in *ai.HealthCheckRequest) (*ai.HealthCheckResponse, error) {
	// todo: add your logic here and delete this line

	return &ai.HealthCheckResponse{}, nil
}
