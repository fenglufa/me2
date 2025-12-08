package dialogue

import (
	"context"

	"github.com/me2/dialogue/rpc/dialogue_client"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话列表
func NewGetSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionsLogic {
	return &GetSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSessionsLogic) GetSessions(req *types.GetSessionsRequest) (resp *types.GetSessionsResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DialogueRpc.GetSessions(l.ctx, &dialogue_client.GetSessionsRequest{
		UserId:   userID,
		AvatarId: req.AvatarId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	sessions := make([]types.SessionInfo, 0, len(rpcResp.Sessions))
	for _, s := range rpcResp.Sessions {
		sessions = append(sessions, types.SessionInfo{
			Id:          s.Id,
			UserId:      s.UserId,
			AvatarId:    s.AvatarId,
			Title:       s.Title,
			LastMessage: s.LastMessage,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	return &types.GetSessionsResponse{
		Sessions: sessions,
		Total:    rpcResp.Total,
	}, nil
}
