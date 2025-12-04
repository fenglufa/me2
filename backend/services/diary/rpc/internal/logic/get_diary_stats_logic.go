package logic

import (
	"context"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiaryStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDiaryStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiaryStatsLogic {
	return &GetDiaryStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取日记统计
func (l *GetDiaryStatsLogic) GetDiaryStats(in *diary.GetDiaryStatsRequest) (*diary.GetDiaryStatsResponse, error) {
	stats, err := l.svcCtx.DiaryModel.GetDiaryStats(l.ctx, in.AvatarId)
	if err != nil {
		l.Errorf("获取日记统计失败: %v", err)
		return nil, err
	}

	return &diary.GetDiaryStatsResponse{
		TotalDiaries:    stats["total_diaries"].(int64),
		AvatarDiaries:   stats["avatar_diaries"].(int64),
		UserDiaries:     stats["user_diaries"].(int64),
		ConsecutiveDays: 0,
		TotalWords:      stats["total_words"].(int64),
		FirstDiaryDate:  stats["first_diary_date"].(string),
		LastDiaryDate:   stats["last_diary_date"].(string),
	}, nil
}
