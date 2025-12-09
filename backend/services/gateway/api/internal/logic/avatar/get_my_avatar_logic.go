package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyAvatarLogic {
	return &GetMyAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyAvatarLogic) GetMyAvatar() (resp *types.AvatarResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.AvatarRpc.GetMyAvatar(l.ctx, &avatar.GetMyAvatarRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	if !rpcResp.HasAvatar {
		return nil, nil
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
