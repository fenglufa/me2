package logic

import (
	"context"

	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarScheduleStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarScheduleStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarScheduleStatusLogic {
	return &GetAvatarScheduleStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分身调度状态
func (l *GetAvatarScheduleStatusLogic) GetAvatarScheduleStatus(in *scheduler.GetAvatarScheduleStatusRequest) (*scheduler.GetAvatarScheduleStatusResponse, error) {
	// 查找调度配置
	existingSchedule, err := l.svcCtx.AvatarScheduleModel.FindByAvatarId(l.ctx, in.AvatarId)
	if err != nil {
		logx.Errorf("Failed to find schedule: %v", err)
		return &scheduler.GetAvatarScheduleStatusResponse{
			ScheduleInfo: nil,
		}, nil
	}

	// 使用 ModelToProtoScheduleInfo 转换并返回
	return &scheduler.GetAvatarScheduleStatusResponse{
		ScheduleInfo: ModelToProtoScheduleInfo(existingSchedule),
	}, nil
}
