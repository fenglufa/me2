package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// 阿里云短信配置
	Aliyun struct {
		AccessKeyId     string
		AccessKeySecret string
		SignName        string // 签名
		TemplateCode    string // 模板代码
	}

	// Redis 配置
	RedisConf redis.RedisConf

	// 短信限制配置
	SmsLimit struct {
		CodeExpire     int // 验证码过期时间（秒），默认 300（5分钟）
		SendInterval   int // 发送间隔时间（秒），默认 60（1分钟）
		DailyLimit     int // 每日发送次数限制，默认 10
	}
}
