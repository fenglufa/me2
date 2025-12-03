package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/me2/scheduler/rpc/internal/model"
	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableAvatarScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnableAvatarScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableAvatarScheduleLogic {
	return &EnableAvatarScheduleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 启用分身调度
func (l *EnableAvatarScheduleLogic) EnableAvatarSchedule(in *scheduler.EnableAvatarScheduleRequest) (*scheduler.EnableAvatarScheduleResponse, error) {
	// 查找现有调度配置
	existingSchedule, err := l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)

	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("Failed to find schedule: %v", err)
		return &scheduler.EnableAvatarScheduleResponse{
			Success: false,
			Message: "Failed to find schedule configuration",
		}, nil
	}

	// 如果不存在，创建新的调度配置
	if err == sql.ErrNoRows {
		newSchedule := &model.AvatarSchedule{
			AvatarId:         in.AvatarId,
			Status:           "active",
			NextScheduleTime: time.Now(),
		}

		_, err := l.svcCtx.AvatarScheduleModel.Insert(l.ctx, newSchedule)
		if err != nil {
			logx.Errorf("Failed to insert schedule: %v", err)
			return &scheduler.EnableAvatarScheduleResponse{
				Success: false,
				Message: "Failed to create schedule configuration",
			}, nil
		}

		// 重新查询以获取完整信息
		existingSchedule, _ = l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)
	} else if existingSchedule.Status == "disabled" {
		// 如果存在且为 disabled，更新为 active
		existingSchedule.Status = "active"
		existingSchedule.NextScheduleTime = time.Now()

		err := l.svcCtx.AvatarScheduleModel.Update(l.ctx, existingSchedule)
		if err != nil {
			logx.Errorf("Failed to update schedule: %v", err)
			return &scheduler.EnableAvatarScheduleResponse{
				Success: false,
				Message: "Failed to enable schedule",
			}, nil
		}
	}

	return &scheduler.EnableAvatarScheduleResponse{
		Success:      true,
		Message:      "Schedule enabled successfully",
		ScheduleInfo: ModelToProtoScheduleInfo(existingSchedule),
	}, nil
}
