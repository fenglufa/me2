package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	avatarID int64
}

func NewUpdateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext, avatarID int64) *UpdateAvatarLogic {
	return &UpdateAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		avatarID: avatarID,
	}
}

func (l *UpdateAvatarLogic) UpdateAvatar(req *types.UpdateAvatarRequest) (resp *types.AvatarResponse, err error) {
	_, err = l.svcCtx.AvatarRpc.UpdateAvatarProfile(l.ctx, &avatar.UpdateAvatarProfileRequest{
		AvatarId:      l.avatarID,
		Nickname:      req.Name,
		AvatarUrl:     req.AvatarUrl,
		Gender:        int32(req.Gender),
		BirthDate:     req.BirthDate,
		Occupation:    req.Occupation,
		MaritalStatus: int32(req.MaritalStatus),
	})
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AvatarRpc.GetAvatarInfo(l.ctx, &avatar.GetAvatarInfoRequest{
		AvatarId: l.avatarID,
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarResponse{
		Id:            rpcResp.Avatar.AvatarId,
		UserId:        rpcResp.Avatar.UserId,
		Name:          rpcResp.Avatar.Nickname,
		AvatarUrl:     rpcResp.Avatar.AvatarUrl,
		Gender:        int64(rpcResp.Avatar.Gender),
		BirthDate:     rpcResp.Avatar.BirthDate,
		Occupation:    rpcResp.Avatar.Occupation,
		MaritalStatus: int64(rpcResp.Avatar.MaritalStatus),
		CreatedAt:     rpcResp.Avatar.CreatedAt,
	}, nil
}
