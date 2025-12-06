package world

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSceneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取场景详情
func NewGetSceneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSceneLogic {
	return &GetSceneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSceneLogic) GetScene(req *types.GetSceneRequest) (resp *types.WorldSceneResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
