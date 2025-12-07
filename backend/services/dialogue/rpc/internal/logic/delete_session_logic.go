package logic

import (
	"context"
	"fmt"

	"github.com/me2/dialogue/rpc/dialogue"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSessionLogic) DeleteSession(in *dialogue.DeleteSessionRequest) (*dialogue.DeleteSessionResponse, error) {
	session, err := l.svcCtx.SessionModel.FindOne(in.SessionId)
	if err != nil {
		return nil, err
	}

	if session.UserId != in.UserId {
		return nil, fmt.Errorf("permission denied")
	}

	if err := l.svcCtx.MessageModel.DeleteBySession(in.SessionId); err != nil {
		return nil, err
	}

	if err := l.svcCtx.SessionModel.Delete(in.SessionId); err != nil {
		return nil, err
	}

	return &dialogue.DeleteSessionResponse{
		Success: true,
	}, nil
}
