package logic

import (
	"context"
	"fmt"

	"github.com/me2/dialogue/rpc/dialogue"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionInfoLogic {
	return &GetSessionInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionInfoLogic) GetSessionInfo(in *dialogue.GetSessionInfoRequest) (*dialogue.GetSessionInfoResponse, error) {
	session, err := l.svcCtx.SessionModel.FindOne(in.SessionId)
	if err != nil {
		return nil, err
	}

	if session.UserId != in.UserId {
		return nil, fmt.Errorf("permission denied")
	}

	return &dialogue.GetSessionInfoResponse{
		Session: &dialogue.SessionInfo{
			Id:          session.Id,
			UserId:      session.UserId,
			AvatarId:    session.AvatarId,
			Title:       session.Title,
			LastMessage: session.LastMessage,
			CreatedAt:   session.CreatedAt,
			UpdatedAt:   session.UpdatedAt,
		},
	}, nil
}
