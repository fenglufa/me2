package logic

import (
	"context"
	"fmt"

	"github.com/me2/sms/rpc/sms"
	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendVerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyCodeLogic {
	return &SendVerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发送验证码
func (l *SendVerifyCodeLogic) SendVerifyCode(in *user.SendVerifyCodeRequest) (*user.SendVerifyCodeResponse, error) {
	if in.Phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}

	// 调用 SMS 服务发送验证码
	_, err := l.svcCtx.SmsRpc.SendCode(l.ctx, &sms.SendCodeRequest{
		Phone: in.Phone,
	})
	if err != nil {
		l.Errorf("发送验证码失败: %v", err)
		return nil, fmt.Errorf("发送验证码失败")
	}

	return &user.SendVerifyCodeResponse{
		Success: true,
	}, nil
}
