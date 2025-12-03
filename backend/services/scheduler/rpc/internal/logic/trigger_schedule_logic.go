package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/me2/action/rpc/action"
	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerScheduleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTriggerScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerScheduleLogic {
	return &TriggerScheduleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 手动触发一次调度
func (l *TriggerScheduleLogic) TriggerSchedule(in *scheduler.TriggerScheduleRequest) (*scheduler.TriggerScheduleResponse, error) {
	// 查找调度配置
	existingSchedule, err := l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)
	if err != nil {
		logx.Errorf("Failed to find schedule: %v", err)
		return &scheduler.TriggerScheduleResponse{
			Success: false,
			Message: "Schedule configuration not found",
		}, nil
	}

	// 调用 Action Service 的 ScheduleAction
	actionResp, err := l.svcCtx.ActionRpc.ScheduleAction(l.ctx, &action.ScheduleActionRequest{
		AvatarId: in.AvatarId,
	})
	if err != nil {
		logx.Errorf("Failed to schedule action: %v", err)
		// 增加失败计数
		existingSchedule.FailedCount++
		l.svcCtx.AvatarScheduleModel.Update(l.ctx, existingSchedule)

		return &scheduler.TriggerScheduleResponse{
			Success: false,
			Message: "Failed to trigger schedule action",
		}, nil
	}

	// 更新调度信息
	existingSchedule.LastScheduleTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	existingSchedule.LastActionType = sql.NullString{
		String: actionResp.Action.ActionType,
		Valid:  true,
	}
	existingSchedule.ScheduleCount++
	existingSchedule.NextScheduleTime = time.Now().Add(1 * time.Hour) // 下次调度时间，可根据业务需求调整

	err = l.svcCtx.AvatarScheduleModel.Update(l.ctx, existingSchedule)
	if err != nil {
		logx.Errorf("Failed to update schedule: %v", err)
	}

	return &scheduler.TriggerScheduleResponse{
		Success:      true,
		Message:      "Schedule triggered successfully",
		ActionType:   actionResp.Action.ActionType,
		ActionLogId:  actionResp.Action.Id,
	}, nil
}
