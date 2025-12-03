package logic

import (
	"context"
	"fmt"

	"github.com/me2/action/rpc/action"
	"github.com/me2/action/rpc/internal/model"
	"github.com/me2/action/rpc/internal/svc"
	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/world/rpc/world_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScheduleActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewScheduleActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleActionLogic {
	return &ScheduleActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 调度行动 - 由定时任务调用，触发分身行动
func (l *ScheduleActionLogic) ScheduleAction(in *action.ScheduleActionRequest) (*action.ScheduleActionResponse, error) {
	// 1. 获取分身信息
	avatarResp, err := l.svcCtx.AvatarRpc.GetAvatarInfo(l.ctx, &avatar_client.GetAvatarInfoRequest{
		AvatarId: in.AvatarId,
	})
	if err != nil {
		return nil, err
	}

	// 2. 计算意图得分
	calculator := NewIntentCalculator()
	intents := calculator.Calculate(avatarResp.Avatar)

	// 3. 选择最佳行为
	bestIntent := calculator.SelectBestAction(intents)

	// 4. 根据行为类型推荐场景
	scenesResp, err := l.svcCtx.WorldRpc.GetScenesForAction(l.ctx, &world_client.GetScenesForActionRequest{
		ActionType: bestIntent.ActionType,
		Limit:      5,
	})
	if err != nil {
		return nil, err
	}

	if len(scenesResp.Recommendations) == 0 {
		return nil, fmt.Errorf("没有找到适合的场景")
	}

	// 5. 选择第一个推荐场景
	selectedScene := scenesResp.Recommendations[0]

	// 6. 记录行动日志
	actionLog := &model.ActionLog{
		AvatarId:      in.AvatarId,
		ActionType:    bestIntent.ActionType,
		SceneId:       selectedScene.Scene.Id,
		SceneName:     selectedScene.Scene.Name,
		IntentScore:   bestIntent.Score,
		TriggerReason: bestIntent.Reason,
		EventId:       0, // 暂时为0，后续由 Event Service 更新
	}

	result, err := l.svcCtx.ActionLogModel.Insert(actionLog)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	actionLog.Id = id

	return &action.ScheduleActionResponse{
		Action: ModelToProtoActionLog(actionLog),
	}, nil
}
