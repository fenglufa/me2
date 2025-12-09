package avatar

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAvatarLogic {
	return &CreateAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAvatarLogic) CreateAvatar(req *types.CreateAvatarRequest) (resp *types.AvatarResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.AvatarRpc.CreateAvatar(l.ctx, &avatar.CreateAvatarRequest{
		UserId:        userID,
		Nickname:      req.Name,
		AvatarUrl:     req.AvatarUrl,
		Gender:        int32(req.Gender),
		BirthDate:     req.BirthDate,
		Occupation:    req.Occupation,
		MaritalStatus: int32(req.MaritalStatus),
	})
	if err != nil {
		return nil, err
	}

	return &types.AvatarResponse{
		Id:            rpcResp.AvatarId,
		UserId:        userID,
		Name:          req.Name,
		AvatarUrl:     req.AvatarUrl,
		Gender:        req.Gender,
		BirthDate:     req.BirthDate,
		Occupation:    req.Occupation,
		MaritalStatus: req.MaritalStatus,
		CreatedAt:     0,
	}, nil
}
