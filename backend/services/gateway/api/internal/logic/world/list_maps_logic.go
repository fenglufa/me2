package world

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/world/rpc/world"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMapsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMapsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMapsLogic {
	return &ListMapsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMapsLogic) ListMaps(req *types.MapListRequest) (resp *types.MapListResponse, err error) {
	rpcResp, err := l.svcCtx.WorldRpc.ListMaps(l.ctx, &world.ListMapsRequest{
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		OnlyActive: req.OnlyActive,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.WorldMapResponse, 0, len(rpcResp.Maps))
	for _, m := range rpcResp.Maps {
		list = append(list, types.WorldMapResponse{
			Id:          m.Id,
			Name:        m.Name,
			Description: m.Description,
			CoverImage:  m.CoverImage,
			IsActive:    m.IsActive,
			CreatedAt:   m.CreatedAt,
		})
	}

	return &types.MapListResponse{
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
