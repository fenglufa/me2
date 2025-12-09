package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新笔记
func NewUpdateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoteLogic {
	return &UpdateNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNoteLogic) UpdateNote(req *types.UpdateNoteRequest) (resp *types.UpdateNoteResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.UpdateNote(l.ctx, &note.UpdateNoteRequest{
		UserId:  userID,
		NoteId:  req.Id,
		RawText: req.RawText,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateNoteResponse{
		Success:   rpcResp.Success,
		Types:     rpcResp.Types,
		AiSummary: rpcResp.AiSummary,
	}, nil
}
