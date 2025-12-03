package logic

import (
	"context"
	"fmt"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/svc"
	"github.com/me2/oss/rpc/oss"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarUploadTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarUploadTokenLogic {
	return &GetAvatarUploadTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAvatarUploadTokenLogic) GetAvatarUploadToken(in *avatar.GetAvatarUploadTokenRequest) (*avatar.GetAvatarUploadTokenResponse, error) {
	if in.AvatarId == 0 {
		return nil, fmt.Errorf("分身ID不能为空")
	}

	resp, err := l.svcCtx.OssRpc.GetUploadToken(l.ctx, &oss.GetUploadTokenRequest{
		ServiceName: "avatar",
		FileName:    in.FileName,
		UserId:      in.AvatarId,
	})
	if err != nil {
		l.Errorf("获取上传凭证失败: %v", err)
		return nil, fmt.Errorf("获取上传凭证失败")
	}

	return &avatar.GetAvatarUploadTokenResponse{
		Host:          resp.Host,
		AccessKeyId:   resp.Accessid,
		Policy:        resp.Policy,
		Signature:     resp.Signature,
		Key:           resp.Dir,
		Expire:        resp.Expire,
		Domain:        resp.Domain,
		CompleteToken: resp.CompleteToken,
	}, nil
}
