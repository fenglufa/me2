package logic

import (
	"context"
	"fmt"

	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscriptionLogic {
	return &GetSubscriptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取订阅信息
func (l *GetSubscriptionLogic) GetSubscription(in *user.GetSubscriptionRequest) (*user.GetSubscriptionResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	u, err := l.svcCtx.UserModel.FindById(in.UserId)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	var expireTime int64
	if u.SubscriptionExpireTime.Valid {
		expireTime = u.SubscriptionExpireTime.Time.Unix()
	}

	return &user.GetSubscriptionResponse{
		SubscriptionType:       u.SubscriptionType,
		SubscriptionExpireTime: expireTime,
	}, nil
}
