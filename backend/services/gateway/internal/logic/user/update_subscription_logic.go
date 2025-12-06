package user

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新订阅
func NewUpdateSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubscriptionLogic {
	return &UpdateSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSubscriptionLogic) UpdateSubscription(req *types.UpdateSubscriptionRequest) (resp *types.SubscriptionResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
