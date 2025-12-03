package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// 阿里云 OSS 配置
	Aliyun struct {
		AccessKeyId     string
		AccessKeySecret string
		Endpoint        string // OSS Endpoint
		Region          string // OSS Region，如 oss-cn-hangzhou
	}

	// JWT 密钥（用于生成 complete_token）
	JwtSecret string

	// 各业务服务的 OSS 配置
	Services map[string]ServiceConfig
}

// ServiceConfig 单个服务的 OSS 配置
type ServiceConfig struct {
	Bucket      string   // Bucket 名称
	Directory   string   // 文件目录前缀
	Domain      string   // 访问域名
	ExpireTime  int64    // 签名过期时间（秒）
	MaxFileSize int64    // 最大文件大小（字节）
	AllowedExts []string // 允许的文件扩展名
}
