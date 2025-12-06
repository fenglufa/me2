package avatar

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarUploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分身头像上传凭证
func NewGetAvatarUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarUploadTokenLogic {
	return &GetAvatarUploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarUploadTokenLogic) GetAvatarUploadToken() (resp *types.AvatarUploadTokenResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
