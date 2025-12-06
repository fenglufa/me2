package avatar

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteAvatarUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 完成分身头像上传
func NewCompleteAvatarUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteAvatarUploadLogic {
	return &CompleteAvatarUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteAvatarUploadLogic) CompleteAvatarUpload(req *types.AvatarUploadCompleteRequest) (resp *types.AvatarUploadCompleteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
