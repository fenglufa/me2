package logic

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/action/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastActionByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastActionByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastActionByUserLogic {
	return &GetLastActionByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户最近一次行动
func (l *GetLastActionByUserLogic) GetLastActionByUser(in *action.GetLastActionByUserRequest) (*action.GetLastActionByUserResponse, error) {
	log, err := l.svcCtx.ActionLogModel.FindLastByUserId(in.UserId)
	if err != nil {
		return nil, err
	}

	if log == nil {
		return &action.GetLastActionByUserResponse{}, nil
	}

	return &action.GetLastActionByUserResponse{
		Action: ModelToProtoActionLog(log),
	}, nil
}
