package world

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMapsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取世界地图列表
func NewListMapsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMapsLogic {
	return &ListMapsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMapsLogic) ListMaps(req *types.MapListRequest) (resp *types.MapListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
