package dialogue

import (
	"context"

	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSessionLogic) DeleteSession(req *types.DeleteSessionRequest) (resp *types.DeleteSessionResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DialogueRpc.DeleteSession(l.ctx, &dialogue_client.DeleteSessionRequest{
		SessionId: req.Id,
		UserId:    userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteSessionResponse{
		Success: rpcResp.Success,
	}, nil
}
