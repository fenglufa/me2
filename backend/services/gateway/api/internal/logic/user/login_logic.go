package user

import (
	"context"

	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/me2/pkg/utils"
	"github.com/me2/user/rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.LoginOrRegister(l.ctx, &user.LoginOrRegisterRequest{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(rpcResp.UserId, l.svcCtx.Config.Auth.AccessSecret, int(l.svcCtx.Config.Auth.AccessExpire/(24*3600)))
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token:    token,
		UserId:   rpcResp.UserId,
		AvatarId: 0,
	}, nil
}
