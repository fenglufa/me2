package logic

import (
	"context"

	"github.com/me2/action/rpc/action"
	"github.com/me2/action/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetActionHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActionHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActionHistoryLogic {
	return &GetActionHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取行动历史
func (l *GetActionHistoryLogic) GetActionHistory(in *action.GetActionHistoryRequest) (*action.GetActionHistoryResponse, error) {
	// 默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 查询行动历史
	logs, total, err := l.svcCtx.ActionLogModel.FindByAvatarId(in.AvatarId, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为 proto 格式
	actions := make([]*action.ActionLog, 0, len(logs))
	for _, log := range logs {
		actions = append(actions, ModelToProtoActionLog(log))
	}

	return &action.GetActionHistoryResponse{
		Actions: actions,
		Total:   total,
	}, nil
}
