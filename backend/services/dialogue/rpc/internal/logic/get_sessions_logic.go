package logic

import (
	"context"

	"github.com/me2/dialogue/rpc/dialogue"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionsLogic {
	return &GetSessionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionsLogic) GetSessions(in *dialogue.GetSessionsRequest) (*dialogue.GetSessionsResponse, error) {
	sessions, err := l.svcCtx.SessionModel.FindByUserAndAvatar(in.UserId, in.AvatarId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.SessionModel.CountByUserAndAvatar(in.UserId, in.AvatarId)
	if err != nil {
		return nil, err
	}

	var sessionInfos []*dialogue.SessionInfo
	for _, s := range sessions {
		sessionInfos = append(sessionInfos, &dialogue.SessionInfo{
			Id:          s.Id,
			UserId:      s.UserId,
			AvatarId:    s.AvatarId,
			Title:       s.Title,
			LastMessage: s.LastMessage,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	return &dialogue.GetSessionsResponse{
		Sessions: sessionInfos,
		Total:    total,
	}, nil
}
