package logic

import (
	"context"
	"fmt"

	"github.com/me2/sms/rpc/internal/svc"
	"github.com/me2/sms/rpc/sms"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VerifyCode 验证验证码
func (l *VerifyCodeLogic) VerifyCode(in *sms.VerifyCodeRequest) (*sms.VerifyCodeResponse, error) {
	// 1. 参数验证
	if in.Phone == "" || in.Code == "" {
		return &sms.VerifyCodeResponse{
			Valid:   false,
			Message: "手机号和验证码不能为空",
		}, nil
	}

	// 2. 从 Redis 获取验证码
	key := fmt.Sprintf("sms:code:%s", in.Phone)
	storedCode, err := l.svcCtx.Redis.Get(key)
	if err != nil {
		l.Errorf("获取验证码失败: phone=%s, err=%v", in.Phone, err)
		return &sms.VerifyCodeResponse{
			Valid:   false,
			Message: "验证码已过期或不存在",
		}, nil
	}

	// 3. 验证码比对
	if storedCode != in.Code {
		l.Errorf("验证码错误: phone=%s", in.Phone)
		return &sms.VerifyCodeResponse{
			Valid:   false,
			Message: "验证码错误",
		}, nil
	}

	// 4. 验证成功，删除验证码
	_, err = l.svcCtx.Redis.Del(key)
	if err != nil {
		l.Errorf("删除验证码失败: phone=%s, err=%v", in.Phone, err)
	}

	l.Infof("验证码验证成功: phone=%s", in.Phone)
	return &sms.VerifyCodeResponse{
		Valid:   true,
		Message: "验证成功",
	}, nil
}
