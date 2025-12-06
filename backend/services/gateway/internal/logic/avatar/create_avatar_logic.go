package avatar

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建分身
func NewCreateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAvatarLogic {
	return &CreateAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAvatarLogic) CreateAvatar(req *types.CreateAvatarRequest) (resp *types.AvatarResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
