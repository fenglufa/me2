package logic

import (
	"context"
	"fmt"

	"github.com/me2/event/rpc/event"
	"github.com/me2/event/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventTimelineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEventTimelineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventTimelineLogic {
	return &GetEventTimelineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetEventTimeline 获取事件时间线
func (l *GetEventTimelineLogic) GetEventTimeline(in *event.GetEventTimelineRequest) (*event.GetEventTimelineResponse, error) {
	// 参数验证
	if in.AvatarId == 0 {
		return nil, fmt.Errorf("avatar_id 不能为空")
	}

	page := in.Page
	pageSize := in.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 查询事件列表
	events, err := l.svcCtx.EventHistoryModel.FindByAvatarId(in.AvatarId, page, pageSize)
	if err != nil {
		l.Errorf("查询事件历史失败: %v", err)
		return nil, fmt.Errorf("查询事件历史失败")
	}

	// 统计总数
	total, err := l.svcCtx.EventHistoryModel.CountByAvatarId(in.AvatarId)
	if err != nil {
		l.Errorf("统计事件数量失败: %v", err)
		return nil, fmt.Errorf("统计事件数量失败")
	}

	// 转换为响应格式
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
