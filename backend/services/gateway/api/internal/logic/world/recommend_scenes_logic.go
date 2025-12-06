package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendScenesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendScenesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendScenesLogic {
	return &RecommendScenesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendScenesLogic) RecommendScenes(req *types.SceneRecommendationRequest) (resp *types.SceneRecommendationResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.GetScenesForAction(l.ctx, &world.GetScenesForActionRequest{
		ActionType: req.ActionType,
		MapId:      req.MapId,
		RegionId:   req.RegionId,
		Limit:      int32(req.Limit),
	})
	if err != nil {
		return nil, err
	}

	recommendations := make([]types.SceneRecommendationItem, 0, len(rpcResp.Recommendations))
	for _, r := range rpcResp.Recommendations {
		recommendations = append(recommendations, types.SceneRecommendationItem{
			Scene: types.WorldSceneResponse{
				Id:              r.Scene.Id,
				RegionId:        r.Scene.RegionId,
				Name:            r.Scene.Name,
				Description:     r.Scene.Description,
				CoverImage:      r.Scene.CoverImage,
				Atmosphere:      r.Scene.Atmosphere,
				Tags:            r.Scene.Tags,
				SuitableActions: r.Scene.SuitableActions,
				Capacity:        int(r.Scene.Capacity),
				IsActive:        r.Scene.IsActive,
				Features: types.SceneFeaturesResponse{
					HasWifi:      r.Scene.Features.HasWifi,
					HasFood:      r.Scene.Features.HasFood,
					HasSeating:   r.Scene.Features.HasSeating,
					IsIndoor:     r.Scene.Features.IsIndoor,
					IsQuiet:      r.Scene.Features.IsQuiet,
					ComfortLevel: int(r.Scene.Features.ComfortLevel),
					SocialLevel:  int(r.Scene.Features.SocialLevel),
				},
				CreatedAt: r.Scene.CreatedAt,
			},
			MatchScore: float64(r.MatchScore),
			Reason:     r.Reason,
		})
	}

	return &types.SceneRecommendationResponse{
		Recommendations: recommendations,
	}, nil
}
