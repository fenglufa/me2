package diary

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarDiariesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分身日记列表
func NewGetAvatarDiariesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarDiariesLogic {
	return &GetAvatarDiariesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarDiariesLogic) GetAvatarDiaries(req *types.DiaryListRequest) (resp *types.DiaryListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
