package user

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户信息
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.UserInfoResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	_, err = l.svcCtx.UserRpc.UpdateUserInfo(l.ctx, &user.UpdateUserInfoRequest{
		UserId:   userID,
		Nickname: req.Nickname,
		Avatar:   req.AvatarUrl,
	})
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:           rpcResp.UserId,
		Phone:            rpcResp.Phone,
		Nickname:         rpcResp.Nickname,
		AvatarUrl:        rpcResp.Avatar,
		SubscriptionTier: rpcResp.SubscriptionType,
		Status:           rpcResp.Status,
		CreatedAt:        rpcResp.CreatedAt,
	}, nil
}
