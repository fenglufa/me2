package world

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListScenesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取区域内的场景列表
func NewListScenesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListScenesLogic {
	return &ListScenesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListScenesLogic) ListScenes(req *types.SceneListRequest) (resp *types.SceneListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
