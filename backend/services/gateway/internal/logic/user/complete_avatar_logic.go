package user

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 完成头像上传
func NewCompleteAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteAvatarLogic {
	return &CompleteAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteAvatarLogic) CompleteAvatar(req *types.AvatarCompleteRequest) (resp *types.AvatarCompleteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
