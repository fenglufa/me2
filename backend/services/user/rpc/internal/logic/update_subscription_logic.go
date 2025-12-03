package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubscriptionLogic {
	return &UpdateSubscriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新订阅状态
func (l *UpdateSubscriptionLogic) UpdateSubscription(in *user.UpdateSubscriptionRequest) (*user.UpdateSubscriptionResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	expireTime := time.Unix(in.SubscriptionExpireTime, 0)
	err := l.svcCtx.UserModel.UpdateSubscription(in.UserId, in.SubscriptionType, expireTime)
	if err != nil {
		l.Errorf("更新订阅失败: %v", err)
		return nil, fmt.Errorf("更新订阅失败")
	}

	l.Infof("更新订阅成功: user_id=%d, type=%d, expire=%d", in.UserId, in.SubscriptionType, in.SubscriptionExpireTime)

	return &user.UpdateSubscriptionResponse{
		Success: true,
	}, nil
}
