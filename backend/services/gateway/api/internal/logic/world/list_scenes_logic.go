package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListScenesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListScenesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListScenesLogic {
	return &ListScenesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListScenesLogic) ListScenes(req *types.SceneListRequest) (resp *types.SceneListResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.ListScenesInRegion(l.ctx, &world.ListScenesInRegionRequest{
		RegionId:   req.RegionId,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		OnlyActive: req.OnlyActive,
		Tags:       req.Tags,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.WorldSceneResponse, 0, len(rpcResp.Scenes))
	for _, s := range rpcResp.Scenes {
		list = append(list, types.WorldSceneResponse{
			Id:              s.Id,
			RegionId:        s.RegionId,
			Name:            s.Name,
			Description:     s.Description,
			CoverImage:      s.CoverImage,
			Atmosphere:      s.Atmosphere,
			Tags:            s.Tags,
			SuitableActions: s.SuitableActions,
			Capacity:        int(s.Capacity),
			IsActive:        s.IsActive,
			Features: types.SceneFeaturesResponse{
				HasWifi:      s.Features.HasWifi,
				HasFood:      s.Features.HasFood,
				HasSeating:   s.Features.HasSeating,
				IsIndoor:     s.Features.IsIndoor,
				IsQuiet:      s.Features.IsQuiet,
				ComfortLevel: int(s.Features.ComfortLevel),
				SocialLevel:  int(s.Features.SocialLevel),
			},
			CreatedAt: s.CreatedAt,
		})
	}

	return &types.SceneListResponse{
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
