# OSS 服务配置说明

## 概述

OSS 服务已重新设计，采用基于策略（Policy）的上传方式，支持多业务场景配置。客户端直接上传到 OSS，服务端提供上传凭证和验证。

## 接口说明

### 1. GetUploadToken - 获取上传凭证

**作用**: 获取 OSS 上传所需的策略、签名和验证令牌

**请求参数**:
```protobuf
message GetUploadTokenRequest {
  string service_name = 1;  // 服务名称: avatar, banner, post
  string file_name = 2;     // 文件名
  int64 user_id = 3;        // 用户 ID
}
```

**响应参数**:
```protobuf
message GetUploadTokenResponse {
  string host = 1;          // OSS 上传地址
  string accessid = 2;      // AccessKeyId
  string policy = 3;        // 上传策略（Base64）
  string signature = 4;     // 签名
  string dir = 6;           // 文件目录前缀
  int64 expire = 7;         // 过期时间戳
  string domain = 8;        // 文件访问域名
  string complete_token = 9; // 完成上传的验证 Token
}
```

### 2. CompleteUpload - 完成上传

**作用**: 验证上传完成，返回文件访问 URL

**请求参数**:
```protobuf
message CompleteUploadRequest {
  string service_name = 1;  // 服务名称
  string key = 2;           // 文件 Key（上传后的完整路径）
  int64 user_id = 3;        // 用户 ID
  string complete_token = 4; // 验证 Token
}
```

**响应参数**:
```protobuf
message CompleteUploadResponse {
  string url = 1;           // 文件访问 URL
}
```

## 配置说明

### 配置文件结构

```yaml
# 阿里云 OSS 配置
Aliyun:
  AccessKeyId: "your-access-key-id"
  AccessKeySecret: "your-access-key-secret"
  Endpoint: "oss-cn-hangzhou.aliyuncs.com"
  Region: "oss-cn-hangzhou"

# JWT 密钥（用于生成和验证 complete_token）
JwtSecret: "your-jwt-secret-key"

# 各业务服务的 OSS 配置
Services:
  avatar:
    Bucket: "your-bucket"
    Directory: "avatar"
    Domain: "https://your-bucket.oss-cn-hangzhou.aliyuncs.com"
    ExpireTime: 3600       # 签名过期时间（秒）
    MaxFileSize: 5242880   # 最大文件大小（字节）
    AllowedExts:           # 允许的文件扩展名
      - jpg
      - jpeg
      - png
```

### 配置参数说明

#### Aliyun - 阿里云配置
- **AccessKeyId**: 阿里云访问密钥 ID
- **AccessKeySecret**: 阿里云访问密钥 Secret
- **Endpoint**: OSS Endpoint，如 `oss-cn-hangzhou.aliyuncs.com`
- **Region**: OSS 区域，如 `oss-cn-hangzhou`

#### JwtSecret - JWT 密钥
用于生成和验证 complete_token，确保上传完成请求的安全性。

#### Services - 业务服务配置

每个业务服务（如 avatar、banner、post）可以有独立的配置：

- **Bucket**: OSS Bucket 名称
- **Directory**: 文件存储目录前缀
- **Domain**: 文件访问域名
- **ExpireTime**: 上传凭证过期时间（秒）
- **MaxFileSize**: 单个文件最大大小（字节）
- **AllowedExts**: 允许上传的文件扩展名列表

### 预设业务配置

#### 1. avatar - 头像上传
- 最大文件大小: 5MB
- 允许格式: jpg, jpeg, png, gif
- 适用场景: 用户头像

#### 2. banner - 横幅上传
- 最大文件大小: 10MB
- 允许格式: jpg, jpeg, png
- 适用场景: 个人主页横幅、活动横幅

#### 3. post - 帖子内容上传
- 最大文件大小: 100MB
- 允许格式: jpg, jpeg, png, gif, mp4, mov, avi
- 适用场景: 帖子中的图片、视频

## 客户端上传流程

### 步骤 1: 获取上传凭证

调用 `GetUploadToken` 接口:

```go
response, err := ossClient.GetUploadToken(ctx, &oss.GetUploadTokenRequest{
    ServiceName: "avatar",
    FileName:    "profile.jpg",
    UserId:      12345,
})
```

### 步骤 2: 使用凭证上传到 OSS

使用返回的凭证通过 OSS SDK 或表单直接上传:

**使用表单上传示例**:
```html
<form method="post" enctype="multipart/form-data" action="{{.Host}}">
    <input type="hidden" name="OSSAccessKeyId" value="{{.Accessid}}">
    <input type="hidden" name="policy" value="{{.Policy}}">
    <input type="hidden" name="signature" value="{{.Signature}}">
    <input type="hidden" name="key" value="{{.Dir}}${filename}">
    <input type="file" name="file">
    <button type="submit">上传</button>
</form>
```

**使用 JavaScript 上传示例**:
```javascript
const formData = new FormData();
formData.append('OSSAccessKeyId', response.accessid);
formData.append('policy', response.policy);
formData.append('signature', response.signature);
formData.append('key', response.dir + fileName);
formData.append('file', fileBlob);

fetch(response.host, {
    method: 'POST',
    body: formData
});
```

### 步骤 3: 完成上传验证

上传成功后，调用 `CompleteUpload` 接口:

```go
urlResponse, err := ossClient.CompleteUpload(ctx, &oss.CompleteUploadRequest{
    ServiceName:   "avatar",
    Key:           uploadedKey,  // OSS 返回的文件 key
    UserId:        12345,
    CompleteToken: response.CompleteToken,
})

// urlResponse.Url 即为文件访问地址
```

## 安全机制

### 1. 策略签名
- 使用 HMAC-SHA1 算法对上传策略进行签名
- 策略包含文件大小限制、目录限制、过期时间等

### 2. JWT 令牌验证
- GetUploadToken 返回的 complete_token 使用 JWT 加密
- CompleteUpload 时验证令牌的有效性
- 令牌中包含 service_name, dir, user_id 等信息

### 3. 多层验证
- 文件扩展名验证
- 文件大小限制
- 目录路径验证
- 用户 ID 匹配验证
- 服务名称匹配验证

## 与旧版本的区别

| 特性 | 旧版本 | 新版本 |
|------|--------|--------|
| 上传方式 | 预签名 URL | 策略 + 签名 |
| 接口 | GetUploadUrl, GetDownloadUrl, DeleteFile | GetUploadToken, CompleteUpload |
| 配置 | 单一配置 | 多业务配置 |
| 安全验证 | 简单签名 | JWT + 多层验证 |
| 文件限制 | 无 | 按业务配置扩展名和大小 |

## 示例配置

### 开发环境 (oss-dev.yaml)

```yaml
Name: oss-rpc
ListenOn: 0.0.0.0:8002

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: oss.rpc

Aliyun:
  AccessKeyId: "dev-key-id"
  AccessKeySecret: "dev-key-secret"
  Endpoint: "oss-cn-hangzhou.aliyuncs.com"
  Region: "oss-cn-hangzhou"

JwtSecret: "dev-jwt-secret-key-change-in-production"

Services:
  avatar:
    Bucket: "dev-bucket"
    Directory: "avatar"
    Domain: "https://dev-bucket.oss-cn-hangzhou.aliyuncs.com"
    ExpireTime: 3600
    MaxFileSize: 5242880
    AllowedExts:
      - jpg
      - jpeg
      - png
      - gif
```

### 生产环境 (oss.yaml)

```yaml
Name: oss-rpc
ListenOn: 0.0.0.0:8002

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: oss.rpc

Aliyun:
  AccessKeyId: ${ALIYUN_ACCESS_KEY_ID}
  AccessKeySecret: ${ALIYUN_ACCESS_KEY_SECRET}
  Endpoint: ${ALIYUN_OSS_ENDPOINT}
  Region: ${ALIYUN_OSS_REGION}

JwtSecret: ${JWT_SECRET}

Services:
  avatar:
    Bucket: ${ALIYUN_OSS_BUCKET_NAME}
    Directory: "avatar"
    Domain: ${ALIYUN_OSS_DOMAIN}
    ExpireTime: 3600
    MaxFileSize: 5242880
    AllowedExts:
      - jpg
      - jpeg
      - png
      - gif
```

## 添加新业务配置

如需添加新的业务场景（如音频上传），只需在配置文件中添加：

```yaml
Services:
  audio:
    Bucket: "your-bucket"
    Directory: "audio"
    Domain: "https://your-bucket.oss-cn-hangzhou.aliyuncs.com"
    ExpireTime: 3600
    MaxFileSize: 52428800  # 50MB
    AllowedExts:
      - mp3
      - wav
      - m4a
```

无需修改代码即可支持新的业务场景。

## 错误处理

常见错误及解决方案：

1. **"不支持的服务类型"**: 检查配置文件中是否配置了对应的 service_name
2. **"不支持的文件类型"**: 检查文件扩展名是否在 AllowedExts 列表中
3. **"无效的验证令牌"**: complete_token 已过期或被篡改
4. **"用户ID不匹配"**: CompleteUpload 的 user_id 与 GetUploadToken 时不一致
5. **"文件路径不合法"**: 上传的文件路径不在允许的目录下

## 监控建议

建议监控以下指标：

1. **上传成功率**: CompleteUpload 调用次数 / GetUploadToken 调用次数
2. **各业务上传量**: 按 service_name 统计
3. **文件大小分布**: 了解用户上传习惯
4. **错误类型分布**: 及时发现配置或使用问题
