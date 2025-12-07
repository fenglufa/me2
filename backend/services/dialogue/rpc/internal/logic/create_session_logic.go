package logic

import (
	"context"
	"time"

	"github.com/me2/dialogue/rpc/dialogue"
	"github.com/me2/dialogue/rpc/internal/model"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(in *dialogue.CreateSessionRequest) (*dialogue.CreateSessionResponse, error) {
	now := time.Now().Unix()

	session := &model.DialogueSession{
		UserId:      in.UserId,
		AvatarId:    in.AvatarId,
		Title:       in.Title,
		LastMessage: "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := model.ValidateSession(session); err != nil {
		return nil, err
	}

	result, err := l.svcCtx.SessionModel.Insert(session)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	return &dialogue.CreateSessionResponse{
		SessionId: id,
	}, nil
}
