package world

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendScenesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据行为类型推荐场景
func NewRecommendScenesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendScenesLogic {
	return &RecommendScenesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendScenesLogic) RecommendScenes(req *types.SceneRecommendationRequest) (resp *types.SceneRecommendationResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
