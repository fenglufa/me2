package user

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/user/rpc/user"
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
	_, err = l.svcCtx.UserRpc.UpdateSubscription(l.ctx, &user.UpdateSubscriptionRequest{
		UserId:                 req.UserId,
		SubscriptionType:       req.SubscriptionTier,
		SubscriptionExpireTime: req.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.UserRpc.GetSubscription(l.ctx, &user.GetSubscriptionRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &types.SubscriptionResponse{
		UserId:           req.UserId,
		SubscriptionTier: rpcResp.SubscriptionType,
		ExpiresAt:        rpcResp.SubscriptionExpireTime,
		AutoRenew:        false,
	}, nil
}
