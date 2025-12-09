package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchNotesLogic {
	return &SearchNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 搜索笔记（用于对话集成）
func (l *SearchNotesLogic) SearchNotes(in *note.SearchNotesRequest) (*note.SearchNotesResponse, error) {
	// todo: add your logic here and delete this line

	return &note.SearchNotesResponse{}, nil
}
