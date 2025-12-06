package action

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLastActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastActionLogic {
	return &GetLastActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLastActionLogic) GetLastAction() (resp *types.ActionLogResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.ActionRpc.GetLastAction(l.ctx, &action.GetLastActionRequest{
		AvatarId: userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.ActionLogResponse{
		Id:            rpcResp.Action.Id,
		AvatarId:      rpcResp.Action.AvatarId,
		ActionType:    rpcResp.Action.ActionType,
		SceneId:       rpcResp.Action.SceneId,
		SceneName:     rpcResp.Action.SceneName,
		IntentScore:   float64(rpcResp.Action.IntentScore),
		TriggerReason: rpcResp.Action.TriggerReason,
		EventId:       rpcResp.Action.EventId,
		CreatedAt:     rpcResp.Action.CreatedAt,
	}, nil
}
