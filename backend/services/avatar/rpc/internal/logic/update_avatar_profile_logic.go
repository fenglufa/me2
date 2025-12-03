package logic

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAvatarProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAvatarProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAvatarProfileLogic {
	return &UpdateAvatarProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAvatarProfileLogic) UpdateAvatarProfile(in *avatar.UpdateAvatarProfileRequest) (*avatar.UpdateAvatarProfileResponse, error) {
	err := l.svcCtx.AvatarModel.UpdateProfile(in.AvatarId, in.Nickname, in.AvatarUrl)
	if err != nil {
		l.Errorf("更新分身资料失败: %v", err)
		return &avatar.UpdateAvatarProfileResponse{Success: false}, err
	}

	return &avatar.UpdateAvatarProfileResponse{Success: true}, nil
}
