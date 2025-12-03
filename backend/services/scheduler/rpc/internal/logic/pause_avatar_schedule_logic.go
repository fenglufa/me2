package logic

import (
	"context"

	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type PauseAvatarScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPauseAvatarScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PauseAvatarScheduleLogic {
	return &PauseAvatarScheduleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 暂停分身调度
func (l *PauseAvatarScheduleLogic) PauseAvatarSchedule(in *scheduler.PauseAvatarScheduleRequest) (*scheduler.PauseAvatarScheduleResponse, error) {
	// 查找调度配置
	existingSchedule, err := l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)
	if err != nil {
		logx.Errorf("Failed to find schedule: %v", err)
		return &scheduler.PauseAvatarScheduleResponse{
			Success: false,
			Message: "Schedule configuration not found",
		}, nil
	}

	// 更新状态为 paused
	err = l.svcCtx.AvatarScheduleModel.UpdateStatus(l.ctx, existingSchedule.AvatarId, "paused")
	if err != nil {
		logx.Errorf("Failed to update schedule status: %v", err)
		return &scheduler.PauseAvatarScheduleResponse{
			Success: false,
			Message: "Failed to pause schedule",
		}, nil
	}

	return &scheduler.PauseAvatarScheduleResponse{
		Success: true,
		Message: "Schedule paused successfully",
	}, nil
}
