package avatar

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取我的分身
func NewGetMyAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyAvatarLogic {
	return &GetMyAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyAvatarLogic) GetMyAvatar() (resp *types.AvatarResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
