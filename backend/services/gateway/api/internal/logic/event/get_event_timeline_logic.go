package event

import (
	"context"

	"github.com/me2/event/rpc/event"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventTimelineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEventTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventTimelineLogic {
	return &GetEventTimelineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventTimelineLogic) GetEventTimeline(req *types.EventListRequest) (resp *types.EventListResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.EventRpc.GetEventTimeline(l.ctx, &event.GetEventTimelineRequest{
		AvatarId: userID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.EventResponse, 0, len(rpcResp.Events))
	for _, e := range rpcResp.Events {
		list = append(list, types.EventResponse{
			Id:        e.EventId,
			EventType: e.EventType,
			EventText: e.EventText,
			ImageUrl:  e.ImageUrl,
			SceneName: e.SceneName,
		})
	}

	return &types.EventListResponse{
		Total:    int64(rpcResp.Total),
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
