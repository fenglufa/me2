package logic

import (
	"context"
	"fmt"

	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户状态
func (l *UpdateUserStatusLogic) UpdateUserStatus(in *user.UpdateUserStatusRequest) (*user.UpdateUserStatusResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if in.Status < 1 || in.Status > 3 {
		return nil, fmt.Errorf("状态值无效")
	}

	err := l.svcCtx.UserModel.UpdateStatus(in.UserId, in.Status)
	if err != nil {
		l.Errorf("更新用户状态失败: %v", err)
		return nil, fmt.Errorf("更新用户状态失败")
	}

	l.Infof("更新用户状态成功: user_id=%d, status=%d", in.UserId, in.Status)

	return &user.UpdateUserStatusResponse{
		Success: true,
	}, nil
}
