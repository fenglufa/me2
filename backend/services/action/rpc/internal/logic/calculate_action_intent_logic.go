package logic

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/action/rpc/internal/svc"
	"github.com/me2/avatar/rpc/avatar_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateActionIntentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateActionIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateActionIntentLogic {
	return &CalculateActionIntentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 计算行动意图 - 计算 6 种行为的得分
func (l *CalculateActionIntentLogic) CalculateActionIntent(in *action.CalculateActionIntentRequest) (*action.CalculateActionIntentResponse, error) {
	// 获取分身信息
	avatarResp, err := l.svcCtx.AvatarRpc.GetAvatarInfo(l.ctx, &avatar_client.GetAvatarInfoRequest{
		AvatarId: in.AvatarId,
	})
	if err != nil {
		return nil, err
	}

	// 计算意图得分
	calculator := NewIntentCalculator()
	intents := calculator.Calculate(avatarResp.Avatar)

	// 选择最佳行为
	bestIntent := calculator.SelectBestAction(intents)

	return &action.CalculateActionIntentResponse{
		Intents:        intents,
		SelectedAction: bestIntent.ActionType,
	}, nil
}
