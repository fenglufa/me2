package logic

import (
	"context"
	"time"
	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatsLogic {
	return &GetStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStatsLogic) GetStats(in *ai.GetStatsRequest) (*ai.GetStatsResponse, error) {
	startDate := in.StartDate
	endDate := in.EndDate
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	stats, err := l.svcCtx.LogModel.GetStats(in.UserId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &ai.GetStatsResponse{
		TotalCalls:        stats.TotalCalls,
		TotalInputTokens:  stats.TotalInputTokens,
		TotalOutputTokens: stats.TotalOutputTokens,
		TotalCost:         stats.TotalCost,
		SuccessCalls:      stats.SuccessCalls,
		ErrorCalls:        stats.ErrorCalls,
		AvgDurationMs:     int64(stats.AvgDurationMs),
	}, nil
}
