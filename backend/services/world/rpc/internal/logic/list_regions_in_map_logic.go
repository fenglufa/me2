package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionsInMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRegionsInMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsInMapLogic {
	return &ListRegionsInMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取地图的所有区域
func (l *ListRegionsInMapLogic) ListRegionsInMap(in *world.ListRegionsInMapRequest) (*world.ListRegionsInMapResponse, error) {
	// 默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	regions, total, err := l.svcCtx.WorldRegionModel.FindByMapId(in.MapId, page, pageSize, in.OnlyActive)
	if err != nil {
		return nil, err
	}

	// 转换为 Proto 类型
	protoRegions := make([]*world.WorldRegion, 0, len(regions))
	for _, r := range regions {
		protoRegions = append(protoRegions, ModelToProtoRegion(r))
	}

	return &world.ListRegionsInMapResponse{
		Regions: protoRegions,
		Total:   total,
	}, nil
}
