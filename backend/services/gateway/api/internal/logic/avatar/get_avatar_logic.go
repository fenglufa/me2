package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	avatarID int64
}

func NewGetAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext, avatarID int64) *GetAvatarLogic {
	return &GetAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		avatarID: avatarID,
	}
}

func (l *GetAvatarLogic) GetAvatar() (resp *types.AvatarResponse, err error) {
	rpcResp, err := l.svcCtx.AvatarRpc.GetAvatarInfo(l.ctx, &avatar.GetAvatarInfoRequest{
		AvatarId: l.avatarID,
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarResponse{
		Id:                       rpcResp.Avatar.AvatarId,
		UserId:                   rpcResp.Avatar.UserId,
		Name:                     rpcResp.Avatar.Nickname,
		AvatarUrl:                rpcResp.Avatar.AvatarUrl,
		Gender:                   int64(rpcResp.Avatar.Gender),
		BirthDate:                rpcResp.Avatar.BirthDate,
		Occupation:               rpcResp.Avatar.Occupation,
		MaritalStatus:            int64(rpcResp.Avatar.MaritalStatus),
		Warmth:                   float64(rpcResp.Avatar.Personality.Warmth),
		Adventurous:              float64(rpcResp.Avatar.Personality.Adventurous),
		Social:                   float64(rpcResp.Avatar.Personality.Social),
		Creative:                 float64(rpcResp.Avatar.Personality.Creative),
		Calm:                     float64(rpcResp.Avatar.Personality.Calm),
		Energetic:                float64(rpcResp.Avatar.Personality.Energetic),
		CreatedAt:                rpcResp.Avatar.CreatedAt,
		PersonalityType:          rpcResp.Avatar.PersonalityType,
		PersonalityTypeName:      rpcResp.Avatar.PersonalityTypeName,
		PersonalityTypeDescription: rpcResp.Avatar.PersonalityTypeDescription,
	}, nil
}
