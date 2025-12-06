package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRegionsLogic) ListRegions(req *types.RegionListRequest) (resp *types.RegionListResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.ListRegionsInMap(l.ctx, &world.ListRegionsInMapRequest{
		MapId:      req.MapId,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		OnlyActive: req.OnlyActive,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.WorldRegionResponse, 0, len(rpcResp.Regions))
	for _, r := range rpcResp.Regions {
		list = append(list, types.WorldRegionResponse{
			Id:          r.Id,
			MapId:       r.MapId,
			Name:        r.Name,
			Description: r.Description,
			CoverImage:  r.CoverImage,
			Atmosphere:  r.Atmosphere,
			Tags:        r.Tags,
			IsActive:    r.IsActive,
			CreatedAt:   r.CreatedAt,
		})
	}

	return &types.RegionListResponse{
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
