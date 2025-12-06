package user

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/user/rpc/user"
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
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.UserRpc.CompleteAvatarUpload(l.ctx, &user.CompleteAvatarUploadRequest{
		UserId:        userID,
		Key:           req.Key,
		CompleteToken: req.Key,
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarCompleteResponse{
		AvatarUrl: rpcResp.AvatarUrl,
	}, nil
}
