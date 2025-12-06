package diary

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDiariesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户日记列表
func NewGetUserDiariesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDiariesLogic {
	return &GetUserDiariesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDiariesLogic) GetUserDiaries(req *types.DiaryListRequest) (resp *types.DiaryListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
