package event

import (
	"context"
	"time"

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

	rpcResp, err := l.svcCtx.EventRpc.GetUserEventTimeline(l.ctx, &event.GetUserEventTimelineRequest{
		UserId:   userID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.EventResponse, 0, len(rpcResp.Events))
	for _, e := range rpcResp.Events {
		var occurredAt int64
		if t, err := time.Parse("2006-01-02 15:04:05", e.OccurredAt); err == nil {
			occurredAt = t.Unix()
		}

		list = append(list, types.EventResponse{
			Id:         e.EventId,
			AvatarId:   0,
			EventType:  e.EventType,
			EventTitle: e.EventTitle,
			EventText:  e.EventText,
			ImageUrl:   e.ImageUrl,
			SceneId:    0,
			SceneName:  e.SceneName,
			OccurredAt: occurredAt,
		})
	}

	return &types.EventListResponse{
		Total:    int64(rpcResp.Total),
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
