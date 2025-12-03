package logic

import (
	"context"
	"fmt"

	"github.com/me2/oss/rpc/oss"
	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

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

// 获取头像上传凭证
func (l *GetAvatarUploadTokenLogic) GetAvatarUploadToken(in *user.GetAvatarUploadTokenRequest) (*user.GetAvatarUploadTokenResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 调用 OSS 服务获取上传凭证
	resp, err := l.svcCtx.OssRpc.GetUploadToken(l.ctx, &oss.GetUploadTokenRequest{
		ServiceName: "avatar",
		FileName:    in.FileName,
		UserId:      in.UserId,
	})
	if err != nil {
		l.Errorf("获取上传凭证失败: %v", err)
		return nil, fmt.Errorf("获取上传凭证失败")
	}

	return &user.GetAvatarUploadTokenResponse{
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
