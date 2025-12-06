package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
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

	rpcResp, err := l.svcCtx.AvatarRpc.GetAvatarUploadToken(l.ctx, &avatar.GetAvatarUploadTokenRequest{
		AvatarId: myAvatarResp.Avatar.AvatarId,
		FileName: "avatar.jpg",
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarUploadTokenResponse{
		Token:     rpcResp.CompleteToken,
		UploadUrl: rpcResp.Host,
		Key:       rpcResp.Key,
	}, nil
}
