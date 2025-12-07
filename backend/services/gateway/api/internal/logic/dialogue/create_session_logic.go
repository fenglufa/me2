package dialogue

import (
	"context"

	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionRequest) (resp *types.CreateSessionResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DialogueRpc.CreateSession(l.ctx, &dialogue_client.CreateSessionRequest{
		UserId:   userID,
		AvatarId: req.AvatarId,
		Title:    req.Title,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateSessionResponse{
		SessionId: rpcResp.SessionId,
	}, nil
}
