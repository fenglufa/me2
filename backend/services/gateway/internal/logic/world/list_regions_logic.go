package world

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRegionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取地图的区域列表
func NewListRegionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRegionsLogic {
	return &ListRegionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRegionsLogic) ListRegions(req *types.RegionListRequest) (resp *types.RegionListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
