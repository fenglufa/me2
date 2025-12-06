package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSceneLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	sceneID int64
}

func NewGetSceneLogic(ctx context.Context, svcCtx *svc.ServiceContext, sceneID int64) *GetSceneLogic {
	return &GetSceneLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		sceneID: sceneID,
	}
}

func (l *GetSceneLogic) GetScene() (resp *types.WorldSceneResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.GetScene(l.ctx, &world.GetSceneRequest{
		SceneId: l.sceneID,
	})
	if err != nil {
		return nil, err
	}

	return &types.WorldSceneResponse{
		Id:              rpcResp.Scene.Id,
		RegionId:        rpcResp.Scene.RegionId,
		Name:            rpcResp.Scene.Name,
		Description:     rpcResp.Scene.Description,
		CoverImage:      rpcResp.Scene.CoverImage,
		Atmosphere:      rpcResp.Scene.Atmosphere,
		Tags:            rpcResp.Scene.Tags,
		SuitableActions: rpcResp.Scene.SuitableActions,
		Capacity:        int(rpcResp.Scene.Capacity),
		IsActive:        rpcResp.Scene.IsActive,
		Features: types.SceneFeaturesResponse{
			HasWifi:      rpcResp.Scene.Features.HasWifi,
			HasFood:      rpcResp.Scene.Features.HasFood,
			HasSeating:   rpcResp.Scene.Features.HasSeating,
			IsIndoor:     rpcResp.Scene.Features.IsIndoor,
			IsQuiet:      rpcResp.Scene.Features.IsQuiet,
			ComfortLevel: int(rpcResp.Scene.Features.ComfortLevel),
			SocialLevel:  int(rpcResp.Scene.Features.SocialLevel),
		},
		CreatedAt: rpcResp.Scene.CreatedAt,
	}, nil
}
