package logic

import (
	"context"

	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResumeAvatarScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResumeAvatarScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeAvatarScheduleLogic {
	return &ResumeAvatarScheduleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 恢复分身调度
func (l *ResumeAvatarScheduleLogic) ResumeAvatarSchedule(in *scheduler.ResumeAvatarScheduleRequest) (*scheduler.ResumeAvatarScheduleResponse, error) {
	// 查找调度配置
	existingSchedule, err := l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)
	if err != nil {
		logx.Errorf("Failed to find schedule: %v", err)
		return &scheduler.ResumeAvatarScheduleResponse{
			Success: false,
			Message: "Schedule configuration not found",
		}, nil
	}

	// 更新状态为 active
	err = l.svcCtx.AvatarScheduleModel.UpdateStatus(l.ctx, existingSchedule.AvatarId, "active")
	if err != nil {
		logx.Errorf("Failed to update schedule status: %v", err)
		return &scheduler.ResumeAvatarScheduleResponse{
			Success: false,
			Message: "Failed to resume schedule",
		}, nil
	}

	// 重新查询以获取最新信息
	existingSchedule, _ = l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)

	return &scheduler.ResumeAvatarScheduleResponse{
		Success:      true,
		Message:      "Schedule resumed successfully",
		ScheduleInfo: ModelToProtoScheduleInfo(existingSchedule),
	}, nil
}
