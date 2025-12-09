package note

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/note/rpc/note"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchNotesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索笔记
func NewSearchNotesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchNotesLogic {
	return &SearchNotesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchNotesLogic) SearchNotes(req *types.SearchNotesRequest) (resp *types.SearchNotesResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.NoteRpc.SearchNotes(l.ctx, &note.SearchNotesRequest{
		UserId:    userID,
		Query:     req.Query,
		Types:     req.Types,
		DateRange: req.DateRange,
		Limit:     int32(req.Limit),
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

	return &types.SearchNotesResponse{
		List: notes,
	}, nil
}
