package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
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
	userID := l.ctx.Value("user_id").(int64)

	myAvatarResp, err := l.svcCtx.AvatarRpc.GetMyAvatar(l.ctx, &avatar.GetMyAvatarRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	if !myAvatarResp.HasAvatar {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AvatarRpc.CompleteAvatarUpload(l.ctx, &avatar.CompleteAvatarUploadRequest{
		AvatarId:      myAvatarResp.Avatar.AvatarId,
		Key:           req.Key,
		CompleteToken: req.Key,
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarUploadCompleteResponse{
		AvatarUrl: rpcResp.AvatarUrl,
	}, nil
}
