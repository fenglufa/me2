package event

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventTimelineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分身事件时间线
func NewGetEventTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventTimelineLogic {
	return &GetEventTimelineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventTimelineLogic) GetEventTimeline(req *types.EventListRequest) (resp *types.EventListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
