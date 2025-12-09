package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取笔记列表
func NewGetNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotesLogic {
	return &GetNotesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotesLogic) GetNotes(req *types.GetNotesRequest) (resp *types.GetNotesResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.GetNotes(l.ctx, &note.GetNotesRequest{
		UserId:   userID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Types:    req.Types,
	})
	if err != nil {
		return nil, err
	}

	// 转换笔记列表
	notes := make([]types.NoteResponse, 0, len(rpcResp.Notes))
	for _, n := range rpcResp.Notes {
		notes = append(notes, types.NoteResponse{
			Id:        n.Id,
			UserId:    n.UserId,
			AvatarId:  n.AvatarId,
			RawText:   n.RawText,
			AiSummary: n.AiSummary,
			Types:     n.Types,
			EmotionData: types.EmotionData{
				Primary: n.EmotionData.Primary,
				Score:   n.EmotionData.Score,
			},
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		})
	}

	return &types.GetNotesResponse{
		Total: rpcResp.Total,
		List:  notes,
	}, nil
}
