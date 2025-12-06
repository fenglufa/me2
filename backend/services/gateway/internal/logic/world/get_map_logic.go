package world

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取世界地图详情
func NewGetMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMapLogic {
	return &GetMapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMapLogic) GetMap(req *types.GetMapRequest) (resp *types.WorldMapResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
