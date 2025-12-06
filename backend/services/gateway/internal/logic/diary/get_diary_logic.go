package diary

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取日记详情
func NewGetDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiaryLogic {
	return &GetDiaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDiaryLogic) GetDiary(req *types.GetDiaryRequest) (resp *types.DiaryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
