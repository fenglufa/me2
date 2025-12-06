package user

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/user/rpc/user"
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
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.UserRpc.GetSubscription(l.ctx, &user.GetSubscriptionRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.SubscriptionResponse{
		UserId:           userID,
		SubscriptionTier: rpcResp.SubscriptionType,
		ExpiresAt:        rpcResp.SubscriptionExpireTime,
		AutoRenew:        false,
	}, nil
}
