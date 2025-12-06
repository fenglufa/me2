package diary

import (
	"context"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarDiariesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvatarDiariesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarDiariesLogic {
	return &GetAvatarDiariesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarDiariesLogic) GetAvatarDiaries(req *types.DiaryListRequest) (resp *types.DiaryListResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DiaryRpc.GetAvatarDiaryList(l.ctx, &diary.GetAvatarDiaryListRequest{
		AvatarId: userID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.DiaryResponse, 0, len(rpcResp.Diaries))
	for _, d := range rpcResp.Diaries {
		list = append(list, types.DiaryResponse{
			Id:        d.Id,
			AvatarId:  d.AvatarId,
			Type:      d.Type,
			Date:      d.Date,
			Title:     d.Title,
			Content:   d.Content,
			Mood:      d.Mood,
			CreatedAt: 0,
		})
	}

	return &types.DiaryListResponse{
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}, nil
}
