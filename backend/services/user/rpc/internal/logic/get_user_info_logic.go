package logic

import (
	"context"
	"fmt"

	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
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

	return &user.GetUserInfoResponse{
		UserId:                 u.UserId,
		Phone:                  u.Phone,
		Nickname:               u.Nickname,
		Avatar:                 u.Avatar,
		SubscriptionType:       u.SubscriptionType,
		SubscriptionExpireTime: expireTime,
		Status:                 u.Status,
		CreatedAt:              u.CreatedAt.Unix(),
		UpdatedAt:              u.UpdatedAt.Unix(),
	}, nil
}
