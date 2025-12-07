package logic

import (
	"context"

	"github.com/me2/event/rpc/event"
	"github.com/me2/event/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEventTimelineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEventTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEventTimelineLogic {
	return &GetUserEventTimelineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取事件时间线（按用户ID）
func (l *GetUserEventTimelineLogic) GetUserEventTimeline(in *event.GetUserEventTimelineRequest) (*event.GetEventTimelineResponse, error) {
	events, err := l.svcCtx.EventHistoryModel.FindByUserId(in.UserId, in.Page, in.PageSize)
	if err != nil {
		l.Errorf("查询用户事件历史失败: %v", err)
		return nil, err
	}

	total, err := l.svcCtx.EventHistoryModel.CountByUserId(in.UserId)
	if err != nil {
		l.Errorf("统计用户事件数量失败: %v", err)
		return nil, err
	}

	eventInfos := make([]*event.EventInfo, 0, len(events))
	for _, e := range events {
		eventInfos = append(eventInfos, &event.EventInfo{
			EventId:    e.Id,
			EventType:  e.EventType,
			EventTitle: e.EventTitle,
			EventText:  e.EventText,
			ImageUrl:   e.ImageUrl,
			SceneName:  e.SceneName,
			OccurredAt: e.OccurredAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &event.GetEventTimelineResponse{
		Events: eventInfos,
		Total:  total,
	}, nil
}
