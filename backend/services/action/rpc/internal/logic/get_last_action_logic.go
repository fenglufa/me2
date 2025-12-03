package logic

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/action/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastActionLogic {
	return &GetLastActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取最近一次行动
func (l *GetLastActionLogic) GetLastAction(in *action.GetLastActionRequest) (*action.GetLastActionResponse, error) {
	log, err := l.svcCtx.ActionLogModel.FindLastByAvatarId(in.AvatarId)
	if err != nil {
		return nil, err
	}

	var protoLog *action.ActionLog
	if log != nil {
		protoLog = ModelToProtoActionLog(log)
	}

	return &action.GetLastActionResponse{
		Action: protoLog,
	}, nil
}
