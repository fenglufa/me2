package action

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryLogic {
	return &GetHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryLogic) GetHistory(req *types.ActionHistoryRequest) (resp *types.ActionHistoryResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.ActionRpc.GetActionHistory(l.ctx, &action.GetActionHistoryRequest{
		AvatarId: userID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.ActionLogResponse, 0, len(rpcResp.Actions))
	for _, a := range rpcResp.Actions {
		list = append(list, types.ActionLogResponse{
			Id:            a.Id,
			AvatarId:      a.AvatarId,
			ActionType:    a.ActionType,
			SceneId:       a.SceneId,
			SceneName:     a.SceneName,
			IntentScore:   float64(a.IntentScore),
			TriggerReason: a.TriggerReason,
			EventId:       a.EventId,
			CreatedAt:     a.CreatedAt,
		})
	}

	return &types.ActionHistoryResponse{
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
