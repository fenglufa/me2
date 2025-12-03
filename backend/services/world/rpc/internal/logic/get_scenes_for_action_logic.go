package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScenesForActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetScenesForActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScenesForActionLogic {
	return &GetScenesForActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据行为类型推荐场景
func (l *GetScenesForActionLogic) GetScenesForAction(in *world.GetScenesForActionRequest) (*world.GetScenesForActionResponse, error) {
	// 默认限制数量
	limit := in.Limit
	if limit <= 0 {
		limit = 10
	}

	// 查询适合该行为类型的场景
	scenes, err := l.svcCtx.WorldSceneModel.FindByActionType(in.ActionType, in.MapId, in.RegionId, limit)
	if err != nil {
		return nil, err
	}

	// 计算匹配度并生成推荐
	recommendations := make([]*world.SceneRecommendation, 0, len(scenes))
	for _, scene := range scenes {
		score := l.svcCtx.WorldSceneModel.CalculateMatchScore(scene, in.ActionType)
		reason := l.svcCtx.WorldSceneModel.GetRecommendationReason(scene, in.ActionType, score)

		recommendations = append(recommendations, &world.SceneRecommendation{
			Scene:      ModelToProtoScene(scene),
			MatchScore: score,
			Reason:     reason,
		})
	}

	return &world.GetScenesForActionResponse{
		Recommendations: recommendations,
	}, nil
}
