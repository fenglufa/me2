package logic

import (
	"context"
	"fmt"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/svc"
	"github.com/me2/oss/rpc/oss"

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

func (l *CompleteAvatarUploadLogic) CompleteAvatarUpload(in *avatar.CompleteAvatarUploadRequest) (*avatar.CompleteAvatarUploadResponse, error) {
	if in.AvatarId == 0 {
		return nil, fmt.Errorf("分身ID不能为空")
	}

	resp, err := l.svcCtx.OssRpc.CompleteUpload(l.ctx, &oss.CompleteUploadRequest{
		ServiceName:   "avatar",
		Key:           in.Key,
		UserId:        in.AvatarId,
		CompleteToken: in.CompleteToken,
	})
	if err != nil {
		l.Errorf("完成上传失败: %v", err)
		return nil, fmt.Errorf("完成上传失败")
	}

	av, err := l.svcCtx.AvatarModel.FindByAvatarId(in.AvatarId)
	if err != nil {
		l.Errorf("查询分身失败: %v", err)
		return nil, fmt.Errorf("分身不存在")
	}

	err = l.svcCtx.AvatarModel.UpdateProfile(av.AvatarId, av.Nickname, resp.Url)
	if err != nil {
		l.Errorf("更新分身头像失败: %v", err)
		return nil, fmt.Errorf("更新分身头像失败")
	}

	return &avatar.CompleteAvatarUploadResponse{
		Success:   true,
		AvatarUrl: resp.Url,
	}, nil
}
