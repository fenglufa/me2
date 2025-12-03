package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/me2/oss/rpc/internal/oss"
	"github.com/me2/oss/rpc/internal/svc"
	ossProto "github.com/me2/oss/rpc/oss"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteUploadLogic {
	return &CompleteUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 完成上传（验证并返回文件 URL）
func (l *CompleteUploadLogic) CompleteUpload(in *ossProto.CompleteUploadRequest) (*ossProto.CompleteUploadResponse, error) {
	// 1. 参数验证
	if in.ServiceName == "" {
		return nil, fmt.Errorf("service_name 不能为空")
	}
	if in.Key == "" {
		return nil, fmt.Errorf("key 不能为空")
	}
	if in.UserId == 0 {
		return nil, fmt.Errorf("user_id 不能为空")
	}
	if in.CompleteToken == "" {
		return nil, fmt.Errorf("complete_token 不能为空")
	}

	// 2. 验证 complete_token
	claims, err := oss.VerifyCompleteToken(in.CompleteToken, l.svcCtx.Config.JwtSecret)
	if err != nil {
		l.Errorf("验证令牌失败: %v", err)
		return nil, fmt.Errorf("无效的验证令牌")
	}

	// 3. 验证 service_name 是否匹配
	if claims.ServiceName != in.ServiceName {
		l.Errorf("服务名称不匹配: token=%s, request=%s", claims.ServiceName, in.ServiceName)
		return nil, fmt.Errorf("服务名称不匹配")
	}

	// 4. 验证 user_id 是否匹配
	if claims.UserID != in.UserId {
		l.Errorf("用户ID不匹配: token=%d, request=%d", claims.UserID, in.UserId)
		return nil, fmt.Errorf("用户ID不匹配")
	}

	// 5. 验证文件路径是否在允许的目录下
	if !strings.HasPrefix(in.Key, claims.Dir) {
		l.Errorf("文件路径不匹配: key=%s, expected_dir=%s", in.Key, claims.Dir)
		return nil, fmt.Errorf("文件路径不合法")
	}

	// 6. 获取服务配置
	serviceConfig, exists := l.svcCtx.Config.Services[in.ServiceName]
	if !exists {
		l.Errorf("未配置的服务类型: %s", in.ServiceName)
		return nil, fmt.Errorf("不支持的服务类型: %s", in.ServiceName)
	}

	// 7. 构建文件访问 URL
	fileUrl := fmt.Sprintf("%s/%s", strings.TrimRight(serviceConfig.Domain, "/"), in.Key)

	l.Infof("上传完成: service=%s, user_id=%d, key=%s, url=%s",
		in.ServiceName, in.UserId, in.Key, fileUrl)

	return &ossProto.CompleteUploadResponse{
		Url: fileUrl,
	}, nil
}
