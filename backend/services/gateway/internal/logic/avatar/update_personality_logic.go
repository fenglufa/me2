package avatar

import (
	"context"

	"github.com/me2/gateway/internal/svc"
	"github.com/me2/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePersonalityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新分身性格
func NewUpdatePersonalityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePersonalityLogic {
	return &UpdatePersonalityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePersonalityLogic) UpdatePersonality(req *types.UpdatePersonalityRequest) (resp *types.AvatarResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
