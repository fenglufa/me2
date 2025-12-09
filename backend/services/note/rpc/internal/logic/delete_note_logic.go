package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoteLogic {
	return &DeleteNoteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除笔记
func (l *DeleteNoteLogic) DeleteNote(in *note.DeleteNoteRequest) (*note.DeleteNoteResponse, error) {
	// todo: add your logic here and delete this line

	return &note.DeleteNoteResponse{}, nil
}
