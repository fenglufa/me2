package logic

import (
	"context"
	"fmt"

	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息（昵称、头像）
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoRequest) (*user.UpdateUserInfoResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 获取当前用户信息
	existingUser, err := l.svcCtx.UserModel.FindById(in.UserId)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	// 使用现有值作为默认值
	nickname := existingUser.Nickname
	avatar := existingUser.Avatar

	// 如果提供了新值，则更新
	if in.Nickname != "" {
		nickname = in.Nickname
	}
	if in.Avatar != "" {
		avatar = in.Avatar
	}

	// 更新数据库
	err = l.svcCtx.UserModel.UpdateInfo(in.UserId, nickname, avatar)
	if err != nil {
		l.Errorf("更新用户信息失败: %v", err)
		return nil, fmt.Errorf("更新用户信息失败")
	}

	return &user.UpdateUserInfoResponse{
		Success: true,
	}, nil
}
