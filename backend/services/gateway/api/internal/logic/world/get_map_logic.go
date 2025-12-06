package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	mapID  int64
}

func NewGetMapLogic(ctx context.Context, svcCtx *svc.ServiceContext, mapID int64) *GetMapLogic {
	return &GetMapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		mapID:  mapID,
	}
}

func (l *GetMapLogic) GetMap() (resp *types.WorldMapResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.GetMap(l.ctx, &world.GetMapRequest{
		MapId: l.mapID,
	})
	if err != nil {
		return nil, err
	}

	return &types.WorldMapResponse{
		Id:          rpcResp.Map.Id,
		Name:        rpcResp.Map.Name,
		Description: rpcResp.Map.Description,
		CoverImage:  rpcResp.Map.CoverImage,
		IsActive:    rpcResp.Map.IsActive,
		CreatedAt:   rpcResp.Map.CreatedAt,
	}, nil
}
