package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除笔记
func NewDeleteNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoteLogic {
	return &DeleteNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNoteLogic) DeleteNote(req *types.DeleteNoteRequest) (resp *types.DeleteNoteResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.DeleteNote(l.ctx, &note.DeleteNoteRequest{
		UserId: userID,
		NoteId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteNoteResponse{
		Success: rpcResp.Success,
	}, nil
}
