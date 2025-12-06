package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePersonalityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新分身性格
func NewUpdatePersonalityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePersonalityLogic {
	return &UpdatePersonalityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePersonalityLogic) UpdatePersonality(req *types.UpdatePersonalityRequest) (resp *types.AvatarResponse, err error) {
	_, err = l.svcCtx.AvatarRpc.UpdatePersonality(l.ctx, &avatar.UpdatePersonalityRequest{
		AvatarId:           req.Id,
		EventId:            0,
		PersonalityChanges: "",
	})
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.AvatarRpc.GetAvatarInfo(l.ctx, &avatar.GetAvatarInfoRequest{
		AvatarId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarResponse{
		Id:          rpcResp.Avatar.AvatarId,
		UserId:      rpcResp.Avatar.UserId,
		Name:        rpcResp.Avatar.Nickname,
		AvatarUrl:   rpcResp.Avatar.AvatarUrl,
		Warmth:      float64(rpcResp.Avatar.Personality.Warmth),
		Adventurous: float64(rpcResp.Avatar.Personality.Adventurous),
		Social:      float64(rpcResp.Avatar.Personality.Social),
		Creative:    float64(rpcResp.Avatar.Personality.Creative),
		Calm:        float64(rpcResp.Avatar.Personality.Calm),
		Energetic:   float64(rpcResp.Avatar.Personality.Energetic),
		CreatedAt:   rpcResp.Avatar.CreatedAt,
	}, nil
}
