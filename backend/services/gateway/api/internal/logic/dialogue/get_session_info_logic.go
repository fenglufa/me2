package dialogue

import (
	"context"

	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSessionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionInfoLogic {
	return &GetSessionInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSessionInfoLogic) GetSessionInfo(req *types.GetSessionInfoRequest) (resp *types.GetSessionInfoResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DialogueRpc.GetSessionInfo(l.ctx, &dialogue_client.GetSessionInfoRequest{
		SessionId: req.Id,
		UserId:    userID,
	})
	if err != nil {
		return nil, err
	}

	s := rpcResp.Session
	return &types.GetSessionInfoResponse{
		Session: types.SessionInfo{
			Id:          s.Id,
			UserId:      s.UserId,
			AvatarId:    s.AvatarId,
			Title:       s.Title,
			LastMessage: s.LastMessage,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		},
	}, nil
}
