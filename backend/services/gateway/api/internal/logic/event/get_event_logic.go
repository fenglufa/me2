package event

import (
	"context"

	"github.com/me2/event/rpc/event"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	eventID int64
}

func NewGetEventLogic(ctx context.Context, svcCtx *svc.ServiceContext, eventID int64) *GetEventLogic {
	return &GetEventLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		eventID: eventID,
	}
}

func (l *GetEventLogic) GetEvent() (resp *types.EventResponse, err error) {
	rpcResp, err := l.svcCtx.EventRpc.GetEventDetail(l.ctx, &event.GetEventDetailRequest{
		EventId: l.eventID,
	})
	if err != nil {
		return nil, err
	}

	return &types.EventResponse{
		Id:        rpcResp.Event.EventId,
		EventType: rpcResp.Event.EventType,
		EventText: rpcResp.Event.EventText,
		ImageUrl:  rpcResp.Event.ImageUrl,
		SceneName: rpcResp.Event.SceneName,
	}, nil
}
