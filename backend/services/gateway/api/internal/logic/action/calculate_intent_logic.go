package action

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateIntentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCalculateIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateIntentLogic {
	return &CalculateIntentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CalculateIntentLogic) CalculateIntent() (resp *types.ActionIntentResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.ActionRpc.CalculateActionIntent(l.ctx, &action.CalculateActionIntentRequest{
		AvatarId: userID,
	})
	if err != nil {
		return nil, err
	}

	intents := make([]types.ActionIntentItem, 0, len(rpcResp.Intents))
	for _, i := range rpcResp.Intents {
		intents = append(intents, types.ActionIntentItem{
			ActionType: i.ActionType,
			Score:      float64(i.Score),
			Reason:     i.Reason,
		})
	}

	return &types.ActionIntentResponse{
		Intents:        intents,
		SelectedAction: rpcResp.SelectedAction,
	}, nil
}
