package action

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateIntentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 计算分身行动意图
func NewCalculateIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateIntentLogic {
	return &CalculateIntentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CalculateIntentLogic) CalculateIntent() (resp *types.ActionIntentResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
