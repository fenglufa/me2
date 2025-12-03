package logic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/me2/sms/rpc/internal/svc"
	"github.com/me2/sms/rpc/sms"
)

type SendCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCodeLogic {
	return &SendCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendCode 发送验证码
func (l *SendCodeLogic) SendCode(in *sms.SendCodeRequest) (*sms.SendCodeResponse, error) {
	// 1. 参数验证
	if in.Phone == "" {
		return &sms.SendCodeResponse{
			Success: false,
			Message: "手机号不能为空",
		}, nil
	}

	// 获取配置的限制参数（带默认值）
	codeExpire := l.svcCtx.Config.SmsLimit.CodeExpire
	if codeExpire == 0 {
		codeExpire = 300 // 默认5分钟
	}
	sendInterval := l.svcCtx.Config.SmsLimit.SendInterval
	if sendInterval == 0 {
		sendInterval = 60 // 默认60秒
	}
	dailyLimit := l.svcCtx.Config.SmsLimit.DailyLimit
	if dailyLimit == 0 {
		dailyLimit = 10 // 默认10次
	}

	// 2. 检查是否已有未使用的验证码
	codeKey := fmt.Sprintf("sms:code:%s", in.Phone)
	existingCode, err := l.svcCtx.Redis.Get(codeKey)
	if err == nil && existingCode != "" {
		l.Infof("验证码尚未过期: phone=%s", in.Phone)
		return &sms.SendCodeResponse{
			Success: false,
			Message: "验证码已发送，请稍后再试",
		}, nil
	}

	// 3. 检查发送频率限制
	limitKey := fmt.Sprintf("sms:limit:%s", in.Phone)
	exists, err := l.svcCtx.Redis.Exists(limitKey)
	if err != nil {
		l.Errorf("检查发送频率失败: phone=%s, err=%v", in.Phone, err)
	}
	if exists {
		l.Infof("发送过于频繁: phone=%s", in.Phone)
		return &sms.SendCodeResponse{
			Success: false,
			Message: fmt.Sprintf("发送过于频繁，请%d秒后重试", sendInterval),
		}, nil
	}

	// 4. 检查每日发送次数限制
	dailyKey := fmt.Sprintf("sms:daily:%s", in.Phone)
	dailyCount, err := l.svcCtx.Redis.Get(dailyKey)
	if err == nil && dailyCount != "" {
		// 解析发送次数
		var count int
		fmt.Sscanf(dailyCount, "%d", &count)
		if count >= dailyLimit {
			l.Infof("超过每日发送限制: phone=%s, count=%d, limit=%d", in.Phone, count, dailyLimit)
			return &sms.SendCodeResponse{
				Success: false,
				Message: "今日发送次数已达上限",
			}, nil
		}
	}

	// 5. 生成6位验证码
	code := generateCode()

	// 6. 调用阿里云短信服务
	client, err := l.createClient()
	if err != nil {
		l.Errorf("创建阿里云客户端失败: %v", err)
		return &sms.SendCodeResponse{
			Success: false,
			Message: "发送失败",
		}, nil
	}

	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(in.Phone),
		SignName:      tea.String(l.svcCtx.Config.Aliyun.SignName),
		TemplateCode:  tea.String(l.svcCtx.Config.Aliyun.TemplateCode),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%s"}`, code)),
	}

	resp, err := client.SendSms(request)
	if err != nil {
		l.Errorf("发送短信失败: phone=%s, err=%v", in.Phone, err)
		return &sms.SendCodeResponse{
			Success: false,
			Message: "发送失败",
		}, nil
	}

	if *resp.Body.Code != "OK" {
		l.Errorf("短信发送失败: phone=%s, code=%s, message=%s",
			in.Phone, *resp.Body.Code, *resp.Body.Message)
		return &sms.SendCodeResponse{
			Success: false,
			Message: "发送失败",
		}, nil
	}

	// 7. 存储验证码到 Redis
	err = l.svcCtx.Redis.Setex(codeKey, code, codeExpire)
	if err != nil {
		l.Errorf("存储验证码失败: phone=%s, err=%v", in.Phone, err)
		return &sms.SendCodeResponse{
			Success: false,
			Message: "发送失败",
		}, nil
	}

	// 8. 设置发送频率限制
	err = l.svcCtx.Redis.Setex(limitKey, "1", sendInterval)
	if err != nil {
		l.Errorf("设置频率限制失败: phone=%s, err=%v", in.Phone, err)
	}

	// 9. 增加每日发送次数
	dailyCount, _ = l.svcCtx.Redis.Get(dailyKey)
	var count int
	if dailyCount != "" {
		fmt.Sscanf(dailyCount, "%d", &count)
	}
	count++

	// 计算到今天结束的秒数
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	ttl := int(endOfDay.Sub(now).Seconds())

	err = l.svcCtx.Redis.Setex(dailyKey, fmt.Sprintf("%d", count), ttl)
	if err != nil {
		l.Errorf("更新每日计数失败: phone=%s, err=%v", in.Phone, err)
	}

	l.Infof("验证码发送成功: phone=%s, daily_count=%d/%d, code_expire=%ds, send_interval=%ds",
		in.Phone, count, dailyLimit, codeExpire, sendInterval)
	return &sms.SendCodeResponse{
		Success: true,
		Message: "发送成功",
	}, nil
}

// createClient 创建阿里云客户端
func (l *SendCodeLogic) createClient() (*dysmsapi.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(l.svcCtx.Config.Aliyun.AccessKeyId),
		AccessKeySecret: tea.String(l.svcCtx.Config.Aliyun.AccessKeySecret),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	return dysmsapi.NewClient(config)
}

// generateCode 生成6位随机验证码
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
