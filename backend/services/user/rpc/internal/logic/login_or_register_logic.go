package logic

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/me2/user/rpc/internal/model"
	"github.com/me2/user/rpc/internal/svc"
	"github.com/me2/user/rpc/user"

	"github.com/me2/sms/rpc/sms"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginOrRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginOrRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginOrRegisterLogic {
	return &LoginOrRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登录或注册（手机号+验证码）
func (l *LoginOrRegisterLogic) LoginOrRegister(in *user.LoginOrRegisterRequest) (*user.LoginOrRegisterResponse, error) {
	// 1. 验证参数
	if in.Phone == "" {
		return nil, fmt.Errorf("手机号不能为空")
	}
	if in.Code == "" {
		return nil, fmt.Errorf("验证码不能为空")
	}

	// 2. 调用 SMS 服务验证验证码
	verifyResp, err := l.svcCtx.SmsRpc.VerifyCode(l.ctx, &sms.VerifyCodeRequest{
		Phone: in.Phone,
		Code:  in.Code,
	})
	if err != nil {
		l.Errorf("验证码验证失败: %v", err)
		return nil, fmt.Errorf("验证码验证失败")
	}
	if !verifyResp.Valid {
		return nil, fmt.Errorf("验证码无效或已过期")
	}

	// 3. 查询用户是否存在
	existingUser, err := l.svcCtx.UserModel.FindByPhone(in.Phone)
	if err != nil && err != sql.ErrNoRows {
		l.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("查询用户失败")
	}

	// 4. 用户存在，直接登录
	if existingUser != nil {
		l.Infof("用户登录成功: user_id=%d, phone=%s", existingUser.UserId, existingUser.Phone)

		var expireTime int64
		if existingUser.SubscriptionExpireTime.Valid {
			expireTime = existingUser.SubscriptionExpireTime.Time.Unix()
		}

		return &user.LoginOrRegisterResponse{
			UserId:                 existingUser.UserId,
			Phone:                  existingUser.Phone,
			Nickname:               existingUser.Nickname,
			Avatar:                 existingUser.Avatar,
			SubscriptionType:       existingUser.SubscriptionType,
			SubscriptionExpireTime: expireTime,
			IsNewUser:              false,
		}, nil
	}

	// 5. 用户不存在，注册新用户
	newUserId := l.svcCtx.IDGen.NextID()

	// 生成默认昵称：前缀 + 8位随机字符
	nickname := l.generateDefaultNickname()
	avatar := l.svcCtx.Config.DefaultAvatar

	newUser := &model.User{
		UserId:           newUserId,
		Phone:            in.Phone,
		Nickname:         nickname,
		Avatar:           avatar,
		SubscriptionType: 0, // 默认免费用户
		Status:           1, // 正常状态
	}

	err = l.svcCtx.UserModel.Insert(newUser)
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, fmt.Errorf("创建用户失败")
	}

	l.Infof("用户注册成功: user_id=%d, phone=%s, nickname=%s", newUserId, in.Phone, nickname)

	return &user.LoginOrRegisterResponse{
		UserId:                 newUserId,
		Phone:                  in.Phone,
		Nickname:               nickname,
		Avatar:                 avatar,
		SubscriptionType:       0,
		SubscriptionExpireTime: 0,
		IsNewUser:              true,
	}, nil
}

// generateDefaultNickname 生成默认昵称：前缀 + 8位随机字符（数字+字母）
func (l *LoginOrRegisterLogic) generateDefaultNickname() string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const length = 8

	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}

	return l.svcCtx.Config.NicknamePrefix + string(result)
}
