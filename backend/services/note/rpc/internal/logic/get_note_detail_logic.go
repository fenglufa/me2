package logic

import (
	"context"

	"github.com/me2/note/rpc/internal/svc"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoteDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoteDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoteDetailLogic {
	return &GetNoteDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取笔记详情
func (l *GetNoteDetailLogic) GetNoteDetail(in *note.GetNoteDetailRequest) (*note.GetNoteDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &note.GetNoteDetailResponse{}, nil
}
