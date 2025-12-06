package avatar

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分身详情
func NewGetAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarLogic {
	return &GetAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarLogic) GetAvatar(req *types.GetAvatarRequest) (resp *types.AvatarResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
