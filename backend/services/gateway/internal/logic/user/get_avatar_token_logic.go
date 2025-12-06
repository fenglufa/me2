package user

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取头像上传凭证
func NewGetAvatarTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarTokenLogic {
	return &GetAvatarTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarTokenLogic) GetAvatarToken() (resp *types.AvatarTokenResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
