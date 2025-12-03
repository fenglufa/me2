package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRegionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRegionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRegionLogic {
	return &GetRegionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取区域详情
func (l *GetRegionLogic) GetRegion(in *world.GetRegionRequest) (*world.GetRegionResponse, error) {
	region, err := l.svcCtx.WorldRegionModel.FindOne(in.RegionId)
	if err != nil {
		return nil, err
	}

	return &world.GetRegionResponse{
		Region: ModelToProtoRegion(region),
	}, nil
}
