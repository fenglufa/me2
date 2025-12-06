package event

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取事件详情
func NewGetEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogic {
	return &GetEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventLogic) GetEvent(req *types.GetEventRequest) (resp *types.EventResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
