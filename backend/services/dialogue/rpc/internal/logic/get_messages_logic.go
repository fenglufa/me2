package logic

import (
	"context"
	"fmt"

	"github.com/me2/dialogue/rpc/dialogue"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *dialogue.GetMessagesRequest) (*dialogue.GetMessagesResponse, error) {
	session, err := l.svcCtx.SessionModel.FindOne(in.SessionId)
	if err != nil {
		return nil, err
	}

	if session.UserId != in.UserId {
		return nil, fmt.Errorf("permission denied")
	}

	messages, err := l.svcCtx.MessageModel.FindBySession(in.SessionId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.MessageModel.CountBySession(in.SessionId)
	if err != nil {
		return nil, err
	}

	var msgList []*dialogue.Message
	for _, m := range messages {
		msgList = append(msgList, &dialogue.Message{
			Id:        m.Id,
			SessionId: m.SessionId,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		})
	}

	return &dialogue.GetMessagesResponse{
		Messages: msgList,
		Total:    total,
	}, nil
}
