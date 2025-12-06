package user

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/user/rpc/user"
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
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.UserRpc.GetAvatarUploadToken(l.ctx, &user.GetAvatarUploadTokenRequest{
		UserId:   userID,
		FileName: "avatar.jpg",
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarTokenResponse{
		Token:     rpcResp.CompleteToken,
		UploadUrl: rpcResp.Host,
		Key:       rpcResp.Key,
		Accessid:  rpcResp.AccessKeyId,
		Policy:    rpcResp.Policy,
		Signature: rpcResp.Signature,
		Dir:       rpcResp.Key,
	}, nil
}
