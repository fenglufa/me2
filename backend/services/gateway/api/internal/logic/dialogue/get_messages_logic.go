package dialogue

import (
	"context"

	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessagesLogic) GetMessages(req *types.GetMessagesRequest) (resp *types.GetMessagesResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DialogueRpc.GetMessages(l.ctx, &dialogue_client.GetMessagesRequest{
		SessionId: req.SessionId,
		UserId:    userID,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]types.Message, 0, len(rpcResp.Messages))
	for _, m := range rpcResp.Messages {
		messages = append(messages, types.Message{
			Id:        m.Id,
			SessionId: m.SessionId,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		})
	}

	return &types.GetMessagesResponse{
		Messages: messages,
		Total:    rpcResp.Total,
	}, nil
}
