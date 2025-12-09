package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建笔记
func NewCreateNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoteLogic {
	return &CreateNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNoteLogic) CreateNote(req *types.CreateNoteRequest) (resp *types.CreateNoteResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.CreateNote(l.ctx, &note.CreateNoteRequest{
		UserId:   userID,
		AvatarId: req.AvatarId,
		RawText:  req.RawText,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateNoteResponse{
		Id:        rpcResp.NoteId,
		Types:     rpcResp.Types,
		AiSummary: rpcResp.AiSummary,
		EmotionData: types.EmotionData{
			Primary: rpcResp.EmotionData.Primary,
			Score:   rpcResp.EmotionData.Score,
		},
	}, nil
}
