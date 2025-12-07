package dialogue

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话详情
func NewGetSessionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionInfoLogic {
	return &GetSessionInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSessionInfoLogic) GetSessionInfo() (resp *types.GetSessionInfoResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
