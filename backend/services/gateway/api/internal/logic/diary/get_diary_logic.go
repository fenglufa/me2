package diary

import (
	"context"

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
	return &types.DiaryResponse{
		Id: l.diaryID,
	}, nil
}
