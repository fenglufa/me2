package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotesLogic {
	return &GetNotesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取笔记列表
func (l *GetNotesLogic) GetNotes(in *note.GetNotesRequest) (*note.GetNotesResponse, error) {
	// todo: add your logic here and delete this line

	return &note.GetNotesResponse{}, nil
}
