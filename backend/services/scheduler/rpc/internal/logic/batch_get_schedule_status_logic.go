package logic

import (
	"context"

	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetScheduleStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetScheduleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetScheduleStatusLogic {
	return &BatchGetScheduleStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量获取调度状态
func (l *BatchGetScheduleStatusLogic) BatchGetScheduleStatus(in *scheduler.BatchGetScheduleStatusRequest) (*scheduler.BatchGetScheduleStatusResponse, error) {
	// 批量查询调度配置
	schedules, err := l.svcCtx.AvatarScheduleModel.BatchFindByAvatarIds(l.ctx, in.AvatarIds)
	if err != nil {
		logx.Errorf("Failed to batch find schedules: %v", err)
		return &scheduler.BatchGetScheduleStatusResponse{
			ScheduleInfos: []*scheduler.ScheduleInfo{},
		}, nil
	}

	// 转换并返回
	scheduleInfos := make([]*scheduler.ScheduleInfo, 0, len(schedules))
	for _, s := range schedules {
		scheduleInfos = append(scheduleInfos, ModelToProtoScheduleInfo(s))
	}

	return &scheduler.BatchGetScheduleStatusResponse{
		ScheduleInfos: scheduleInfos,
	}, nil
}
