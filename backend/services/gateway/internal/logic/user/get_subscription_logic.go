package user

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取订阅信息
func NewGetSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscriptionLogic {
	return &GetSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscriptionLogic) GetSubscription() (resp *types.SubscriptionResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
