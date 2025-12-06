package logic

import (
	"context"
	"database/sql"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetDiaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiaryLogic {
	return &GetDiaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDiaryLogic) GetDiary(in *diary.GetDiaryRequest) (*diary.GetDiaryResponse, error) {
	diaryModel, err := l.svcCtx.DiaryModel.FindOne(l.ctx, in.DiaryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "日记不存在")
		}
		return nil, err
	}

	return &diary.GetDiaryResponse{
		Diary: &diary.DiaryInfo{
			Id:           diaryModel.Id,
			AvatarId:     diaryModel.AvatarId,
			Type:         diaryModel.Type,
			Date:         diaryModel.Date.Format("2006-01-02"),
			Title:        diaryModel.Title,
			Content:      diaryModel.Content,
			Mood:         diaryModel.Mood,
			Tags:         []string{},
			ReplyContent: diaryModel.ReplyContent.String,
			EmotionScore: int32(diaryModel.EmotionScore),
			IsImportant:  diaryModel.IsImportant == 1,
			CreatedAt:    diaryModel.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
