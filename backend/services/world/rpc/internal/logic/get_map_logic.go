package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMapLogic {
	return &GetMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取世界地图
func (l *GetMapLogic) GetMap(in *world.GetMapRequest) (*world.GetMapResponse, error) {
	worldMap, err := l.svcCtx.WorldMapModel.FindOne(in.MapId)
	if err != nil {
		return nil, err
	}

	return &world.GetMapResponse{
		Map: ModelToProtoMap(worldMap),
	}, nil
}
