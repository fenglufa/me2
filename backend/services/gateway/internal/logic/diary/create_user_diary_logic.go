package diary

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserDiaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建用户日记
func NewCreateUserDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserDiaryLogic {
	return &CreateUserDiaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserDiaryLogic) CreateUserDiary(req *types.CreateUserDiaryRequest) (resp *types.DiaryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
