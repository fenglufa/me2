package logic

import (
	"context"

	"github.com/me2/world/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSceneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSceneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSceneLogic {
	return &GetSceneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取场景详情
func (l *GetSceneLogic) GetScene(in *world.GetSceneRequest) (*world.GetSceneResponse, error) {
	scene, err := l.svcCtx.WorldSceneModel.FindOne(in.SceneId)
	if err != nil {
		return nil, err
	}

	return &world.GetSceneResponse{
		Scene: ModelToProtoScene(scene),
	}, nil
}
