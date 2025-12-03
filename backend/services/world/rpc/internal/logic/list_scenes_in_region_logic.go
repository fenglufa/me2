package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListScenesInRegionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListScenesInRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListScenesInRegionLogic {
	return &ListScenesInRegionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 列出区域内的场景
func (l *ListScenesInRegionLogic) ListScenesInRegion(in *world.ListScenesInRegionRequest) (*world.ListScenesInRegionResponse, error) {
	// 默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	scenes, total, err := l.svcCtx.WorldSceneModel.FindByRegionId(in.RegionId, page, pageSize, in.OnlyActive, in.Tags)
	if err != nil {
		return nil, err
	}

	// 转换为 Proto 类型
	protoScenes := make([]*world.WorldScene, 0, len(scenes))
	for _, s := range scenes {
		protoScenes = append(protoScenes, ModelToProtoScene(s))
	}

	return &world.ListScenesInRegionResponse{
		Scenes: protoScenes,
		Total:  total,
	}, nil
}
