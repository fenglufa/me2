package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRegionLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	regionID int64
}

func NewGetRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext, regionID int64) *GetRegionLogic {
	return &GetRegionLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		regionID: regionID,
	}
}

func (l *GetRegionLogic) GetRegion() (resp *types.WorldRegionResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.GetRegion(l.ctx, &world.GetRegionRequest{
		RegionId: l.regionID,
	})
	if err != nil {
		return nil, err
	}

	return &types.WorldRegionResponse{
		Id:          rpcResp.Region.Id,
		MapId:       rpcResp.Region.MapId,
		Name:        rpcResp.Region.Name,
		Description: rpcResp.Region.Description,
		CoverImage:  rpcResp.Region.CoverImage,
		Atmosphere:  rpcResp.Region.Atmosphere,
		Tags:        rpcResp.Region.Tags,
		IsActive:    rpcResp.Region.IsActive,
		CreatedAt:   rpcResp.Region.CreatedAt,
	}, nil
}
