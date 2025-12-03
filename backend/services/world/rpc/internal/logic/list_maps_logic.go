package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMapsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMapsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMapsLogic {
	return &ListMapsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取所有地图列表
func (l *ListMapsLogic) ListMaps(in *world.ListMapsRequest) (*world.ListMapsResponse, error) {
	// 默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	maps, total, err := l.svcCtx.WorldMapModel.FindAll(page, pageSize, in.OnlyActive)
	if err != nil {
		return nil, err
	}

	// 转换为 Proto 类型
	protoMaps := make([]*world.WorldMap, 0, len(maps))
	for _, m := range maps {
		protoMaps = append(protoMaps, ModelToProtoMap(m))
	}

	return &world.ListMapsResponse{
		Maps:  protoMaps,
		Total: total,
	}, nil
}
