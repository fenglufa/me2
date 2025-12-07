package dialogue

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话列表
func NewGetSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionsLogic {
	return &GetSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSessionsLogic) GetSessions(req *types.GetSessionsRequest) (resp *types.GetSessionsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
