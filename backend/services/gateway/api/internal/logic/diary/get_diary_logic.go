package diary

import (
	"context"
	"time"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiaryLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	diaryID int64
}

func NewGetDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext, diaryID int64) *GetDiaryLogic {
	return &GetDiaryLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		diaryID: diaryID,
	}
}

func (l *GetDiaryLogic) GetDiary() (resp *types.DiaryResponse, err error) {
	rpcResp, err := l.svcCtx.DiaryRpc.GetDiary(l.ctx, &diary.GetDiaryRequest{
		DiaryId: l.diaryID,
	})
	if err != nil {
		return nil, err
	}

	diaryInfo := rpcResp.Diary
	createdAt, _ := time.Parse("2006-01-02 15:04:05", diaryInfo.CreatedAt)
	return &types.DiaryResponse{
		Id:           diaryInfo.Id,
		AvatarId:     diaryInfo.AvatarId,
		Type:         diaryInfo.Type,
		Date:         diaryInfo.Date,
		Title:        diaryInfo.Title,
		Content:      diaryInfo.Content,
		Mood:         diaryInfo.Mood,
		ReplyContent: diaryInfo.ReplyContent,
		CreatedAt:    createdAt.Unix(),
	}, nil
}
