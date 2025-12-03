package logic

import (
	"context"
	"fmt"

	"github.com/me2/oss/rpc/oss"
	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteAvatarUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteAvatarUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteAvatarUploadLogic {
	return &CompleteAvatarUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 完成头像上传
func (l *CompleteAvatarUploadLogic) CompleteAvatarUpload(in *user.CompleteAvatarUploadRequest) (*user.CompleteAvatarUploadResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 调用 OSS 服务完成上传
	resp, err := l.svcCtx.OssRpc.CompleteUpload(l.ctx, &oss.CompleteUploadRequest{
		ServiceName:   "avatar",
		Key:           in.Key,
		UserId:        in.UserId,
		CompleteToken: in.CompleteToken,
	})
	if err != nil {
		l.Errorf("完成上传失败: %v", err)
		return nil, fmt.Errorf("完成上传失败")
	}

	// 更新用户头像
	u, err := l.svcCtx.UserModel.FindById(in.UserId)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("用户不存在")
	}

	u.Avatar = resp.Url
	err = l.svcCtx.UserModel.UpdateInfo(in.UserId, u.Nickname, u.Avatar)
	if err != nil {
		l.Errorf("更新用户头像失败: %v", err)
		return nil, fmt.Errorf("更新用户头像失败")
	}

	l.Infof("用户头像更新成功: user_id=%d, avatar=%s", in.UserId, resp.Url)

	return &user.CompleteAvatarUploadResponse{
		Success:   true,
		AvatarUrl: resp.Url,
	}, nil
}
