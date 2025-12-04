package logic

import (
	"context"
	"fmt"

	"github.com/me2/event/rpc/event"
	"github.com/me2/event/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEventDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventDetailLogic {
	return &GetEventDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetEventDetail 获取事件详情
func (l *GetEventDetailLogic) GetEventDetail(in *event.GetEventDetailRequest) (*event.GetEventDetailResponse, error) {
	if in.EventId == 0 {
		return nil, fmt.Errorf("event_id 不能为空")
	}

	// 查询事件
	e, err := l.svcCtx.EventHistoryModel.FindOne(in.EventId)
	if err != nil {
		l.Errorf("查询事件失败: %v", err)
		return nil, fmt.Errorf("查询事件失败")
	}

	return &event.GetEventDetailResponse{
		Event: &event.EventInfo{
			EventId:    e.Id,
			EventType:  e.EventType,
			EventTitle: e.EventTitle,
			EventText:  e.EventText,
			ImageUrl:   e.ImageUrl,
			SceneName:  e.SceneName,
			OccurredAt: e.OccurredAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
