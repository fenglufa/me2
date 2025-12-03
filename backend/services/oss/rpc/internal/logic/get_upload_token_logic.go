package logic

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/me2/oss/rpc/internal/oss"
	"github.com/me2/oss/rpc/internal/svc"
	ossProto "github.com/me2/oss/rpc/oss"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadTokenLogic {
	return &GetUploadTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取上传 Token（包含 policy 和 signature）
func (l *GetUploadTokenLogic) GetUploadToken(in *ossProto.GetUploadTokenRequest) (*ossProto.GetUploadTokenResponse, error) {
	// 1. 参数验证
	if in.ServiceName == "" {
		return nil, fmt.Errorf("service_name 不能为空")
	}
	if in.FileName == "" {
		return nil, fmt.Errorf("file_name 不能为空")
	}
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}

	// 2. 获取服务配置
	serviceConfig, exists := l.svcCtx.Config.Services[in.ServiceName]
	if !exists {
		l.Errorf("未配置的服务类型: %s", in.ServiceName)
		return nil, fmt.Errorf("不支持的服务类型: %s", in.ServiceName)
	}

	// 3. 验证文件扩展名
	ext := strings.ToLower(path.Ext(in.FileName))
	if ext != "" && ext[0] == '.' {
		ext = ext[1:] // 去掉开头的点
	}

	validExt := false
	for _, allowedExt := range serviceConfig.AllowedExts {
		if ext == strings.ToLower(allowedExt) {
			validExt = true
			break
		}
	}
	if !validExt {
		l.Errorf("不支持的文件类型: %s, 允许的类型: %v", ext, serviceConfig.AllowedExts)
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 4. 设置文件目录前缀
	dir := serviceConfig.Directory
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	// 5. 生成上传策略
	policy, expireTime, err := oss.GeneratePolicy(dir, serviceConfig.MaxFileSize, serviceConfig.ExpireTime)
	if err != nil {
		l.Errorf("生成策略失败: %v", err)
		return nil, fmt.Errorf("生成上传策略失败")
	}

	// 6. 生成签名
	signature := oss.GenerateSignature(policy, l.svcCtx.Config.Aliyun.AccessKeySecret)

	// 7. 生成完成上传的验证 Token
	completeToken, err := oss.GenerateCompleteToken(
		in.ServiceName,
		dir,
		in.UserId,
		l.svcCtx.Config.JwtSecret,
		serviceConfig.ExpireTime,
	)
	if err != nil {
		l.Errorf("生成完成令牌失败: %v", err)
		return nil, fmt.Errorf("生成验证令牌失败")
	}

	// 8. 构建 OSS Host
	host := fmt.Sprintf("https://%s.%s", serviceConfig.Bucket, l.svcCtx.Config.Aliyun.Endpoint)

	l.Infof("生成上传令牌成功: service=%s, user_id=%d, dir=%s, expire=%d",
		in.ServiceName, in.UserId, dir, expireTime)

	return &ossProto.GetUploadTokenResponse{
		Host:          host,
		Accessid:      l.svcCtx.Config.Aliyun.AccessKeyId,
		Policy:        policy,
		Signature:     signature,
		Dir:           dir,
		Expire:        expireTime,
		Domain:        serviceConfig.Domain,
		CompleteToken: completeToken,
	}, nil
}
