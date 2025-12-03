# User Service

用户服务，提供用户注册、登录、信息管理和订阅管理功能。

## 目录

- [功能](#功能)
- [技术栈](#技术栈)
- [配置](#配置)
- [数据库](#数据库)
- [使用](#使用)
- [API](#api)
  - [SendVerifyCode - 发送验证码](#sendverifycode---发送验证码)
  - [LoginOrRegister - 登录或注册](#loginorregister---登录或注册)
  - [GetUserInfo - 获取用户信息](#getuserinfo---获取用户信息)
  - [UpdateUserInfo - 更新用户信息](#updateuserinfo---更新用户信息)
  - [GetAvatarUploadToken - 获取头像上传凭证](#getavataruploadtoken---获取头像上传凭证)
  - [CompleteAvatarUpload - 完成头像上传](#completeavatarupload---完成头像上传)
  - [UpdateSubscription - 更新订阅状态](#updatesubscription---更新订阅状态)
  - [GetSubscription - 获取订阅信息](#getsubscription---获取订阅信息)
  - [UpdateUserStatus - 更新用户状态](#updateuserstatus---更新用户状态)
- [User ID 生成规则](#user-id-生成规则)
- [默认昵称和头像](#默认昵称和头像)
- [头像上传流程](#头像上传流程)
- [依赖服务](#依赖服务)
- [注意事项](#注意事项)
- [测试](#测试)

## 功能

- 发送验证码 (SendVerifyCode) - 调用 SMS 服务发送验证码
- 登录/注册 (LoginOrRegister) - 手机号+验证码，自动判断新老用户，新用户自动生成随机昵称和默认头像
- 获取用户信息 (GetUserInfo) - 包含昵称、头像等基础信息
- 更新用户信息 (UpdateUserInfo) - 更新昵称、头像
- 获取头像上传凭证 (GetAvatarUploadToken) - 调用 OSS 服务获取上传凭证
- 完成头像上传 (CompleteAvatarUpload) - 完成上传并更新用户头像
- 更新订阅状态 (UpdateSubscription)
- 获取订阅信息 (GetSubscription)
- 更新用户状态 (UpdateUserStatus)

## 技术栈

- go-zero RPC
- MySQL (用户数据存储)
- Etcd (服务注册与发现)
- 雪花算法 (生成 9 位数字 user_id)

## 配置

配置文件: `rpc/etc/user-dev.yaml`

```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8003

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: user.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me2?charset=utf8mb4&parseTime=true&loc=Local

SmsRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: sms.rpc

OssRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: oss.rpc

MachineID: 1  # 雪花算法机器ID (0-1023)

# 默认头像和昵称配置
DefaultAvatar: "https://example.com/default-avatar.png"
NicknamePrefix: "用户"
```

## 数据库

### 创建数据库表

```bash
mysql -u root -p < rpc/user.sql
```

### 表结构

```sql
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键',
    user_id BIGINT UNIQUE NOT NULL COMMENT '用户ID（9位数字）',
    phone VARCHAR(20) UNIQUE NOT NULL COMMENT '手机号',
    nickname VARCHAR(50) DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(500) DEFAULT '' COMMENT '头像URL',
    subscription_type TINYINT DEFAULT 0 COMMENT '订阅类型 0:免费 1:付费',
    subscription_expire_time DATETIME COMMENT '会员过期时间',
    status TINYINT DEFAULT 1 COMMENT '状态 1:正常 2:封禁 3:注销',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 使用

```bash
# 1. 生成 RPC 代码 (必须先执行)
make gen-rpc

# 2. 初始化依赖
make init

# 3. 创建数据库表
mysql -u root -p < rpc/user.sql

# 4. 编译
make build

# 5. 运行（开发模式）
make run-dev

# 6. 测试
make test
```

**注意**: 必须先执行 `make gen-rpc` 生成 proto 代码，再执行 `make init` 初始化依赖。

## API

### SendVerifyCode - 发送验证码

请求:
```protobuf
message SendVerifyCodeRequest {
  string phone = 1;  // 手机号
}
```

响应:
```protobuf
message SendVerifyCodeResponse {
  bool success = 1;
}
```

### LoginOrRegister - 登录或注册

请求:
```protobuf
message LoginOrRegisterRequest {
  string phone = 1;  // 手机号
  string code = 2;   // 验证码
}
```

响应:
```protobuf
message LoginOrRegisterResponse {
  int64 user_id = 1;                    // 用户ID（9位数字）
  string phone = 2;                     // 手机号
  string nickname = 3;                  // 昵称（新用户自动生成）
  string avatar = 4;                    // 头像URL（新用户使用默认头像）
  int32 subscription_type = 5;          // 订阅类型
  int64 subscription_expire_time = 6;   // 会员过期时间
  bool is_new_user = 7;                 // 是否新用户
}
```

### GetUserInfo - 获取用户信息

请求:
```protobuf
message GetUserInfoRequest {
  int64 user_id = 1;
}
```

响应:
```protobuf
message GetUserInfoResponse {
  int64 user_id = 1;
  string phone = 2;
  string nickname = 3;
  string avatar = 4;
  int32 subscription_type = 5;
  int64 subscription_expire_time = 6;
  int32 status = 7;  // 1:正常 2:封禁 3:注销
  int64 created_at = 8;
  int64 updated_at = 9;
}
```

### UpdateUserInfo - 更新用户信息

请求:
```protobuf
message UpdateUserInfoRequest {
  int64 user_id = 1;
  string nickname = 2;  // 昵称（可选）
  string avatar = 3;    // 头像URL（可选）
}
```

响应:
```protobuf
message UpdateUserInfoResponse {
  bool success = 1;
}
```

### GetAvatarUploadToken - 获取头像上传凭证

请求:
```protobuf
message GetAvatarUploadTokenRequest {
  int64 user_id = 1;    // 用户ID
  string file_name = 2; // 文件名
}
```

响应:
```protobuf
message GetAvatarUploadTokenResponse {
  string host = 1;           // OSS 上传地址
  string access_key_id = 2;  // AccessKeyId
  string policy = 3;         // 上传策略（Base64）
  string signature = 4;      // 签名
  string key = 5;            // 文件路径
  int64 expire = 6;          // 过期时间戳
  string domain = 7;         // 文件访问域名
  string complete_token = 8; // 完成上传的验证 Token
}
```

### CompleteAvatarUpload - 完成头像上传

请求:
```protobuf
message CompleteAvatarUploadRequest {
  int64 user_id = 1;         // 用户ID
  string key = 2;            // 文件 Key（上传后的完整路径）
  string complete_token = 3; // 验证 Token
}
```

响应:
```protobuf
message CompleteAvatarUploadResponse {
  bool success = 1;
  string avatar_url = 2; // 头像访问 URL
}
```

### UpdateSubscription - 更新订阅状态

请求:
```protobuf
message UpdateSubscriptionRequest {
  int64 user_id = 1;
  int32 subscription_type = 2;          // 0:免费 1:付费
  int64 subscription_expire_time = 3;   // Unix时间戳
}
```

### GetSubscription - 获取订阅信息

请求:
```protobuf
message GetSubscriptionRequest {
  int64 user_id = 1;
}
```

响应:
```protobuf
message GetSubscriptionResponse {
  int32 subscription_type = 1;
  int64 subscription_expire_time = 2;
}
```

### UpdateUserStatus - 更新用户状态

请求:
```protobuf
message UpdateUserStatusRequest {
  int64 user_id = 1;
  int32 status = 2;  // 1:正常 2:封禁 3:注销
}
```

## User ID 生成规则

使用改进的雪花算法生成 9 位数字的 user_id：

- **标准 Snowflake 结构**：
  - 时间戳（41位）+ Worker ID（10位）+ 随机数（12位）
  - 使用加密级随机数（crypto/rand）

- **混淆函数**：
  - 使用大质数和位运算将 64 位 ID 映射到 9 位数字
  - 确保 ID 看起来完全随机，无规律可循

- **特点**：
  - 范围: 100000000 - 999999999
  - 分布式唯一（通过 Worker ID）
  - 高度随机化（通过混淆函数）
  - 无碰撞（时间戳 + 随机数保证）

## 默认昵称和头像

新用户注册时自动生成：

- **昵称**：`配置前缀 + 8位随机字符`
  - 随机字符包含数字和大小写字母
  - 使用加密级随机数生成器（crypto/rand）
  - 示例：`用户a3Bx9Kp2`

- **头像**：使用配置文件中的默认头像 URL
  - 可在 `user-dev.yaml` 中配置 `DefaultAvatar`
  - 可在 `user-dev.yaml` 中配置 `NicknamePrefix`

## 头像上传流程

用户可以通过调用头像上传接口更新自己的头像。整个流程采用**两阶段上传模式**：

### 1. 获取上传凭证

客户端调用 `GetAvatarUploadToken` 获取上传所需的凭证：

```bash
grpcurl -plaintext -proto user.proto \
  -d '{"user_id": 123456789, "file_name": "avatar.jpg"}' \
  127.0.0.1:8003 user.User/GetAvatarUploadToken
```

返回示例：
```json
{
  "host": "https://bucket.oss-cn-hangzhou.aliyuncs.com",
  "accessKeyId": "your-access-key-id",
  "policy": "base64-encoded-policy",
  "signature": "signature-string",
  "key": "avatar/",
  "expire": "1764750224",
  "domain": "https://bucket.oss-cn-hangzhou.aliyuncs.com",
  "completeToken": "jwt-token"
}
```

### 2. 客户端直接上传到 OSS

客户端使用获取的凭证，直接通过 HTTP POST 上传文件到 OSS：

```bash
curl -X POST "https://bucket.oss-cn-hangzhou.aliyuncs.com" \
  -F "key=avatar/filename.jpg" \
  -F "policy=base64-encoded-policy" \
  -F "OSSAccessKeyId=your-access-key-id" \
  -F "signature=signature-string" \
  -F "file=@/path/to/avatar.jpg"
```

### 3. 完成上传并更新数据库

客户端上传成功后，调用 `CompleteAvatarUpload` 通知服务端：

```bash
grpcurl -plaintext -proto user.proto \
  -d '{
    "user_id": 123456789,
    "key": "avatar/filename.jpg",
    "complete_token": "jwt-token"
  }' \
  127.0.0.1:8003 user.User/CompleteAvatarUpload
```

返回示例：
```json
{
  "success": true,
  "avatarUrl": "https://bucket.oss-cn-hangzhou.aliyuncs.com/avatar/filename.jpg"
}
```

### 流程说明

1. **安全性**：User 服务调用 OSS 服务生成临时上传凭证，凭证包含过期时间和权限限制
2. **性能**：客户端直接上传到 OSS，不经过应用服务器，减少服务器带宽压力
3. **验证**：完成上传时通过 `complete_token` 验证，确保上传的合法性
4. **原子性**：只有在 OSS 验证成功后，才会更新数据库中的头像 URL

## 依赖服务

- **SMS Service**: 验证码验证
- **OSS Service**: 头像上传凭证和完成上传
- **MySQL**: 用户数据存储
- **Etcd**: 服务注册与发现

## 注意事项

- User 服务存储基础用户信息（昵称、头像），Avatar 服务存储 AI 生成的虚拟形象
- 无密码设计，纯手机号 + 验证码登录
- JWT 在 Gateway 层生成，User Service 只负责验证用户身份
- 登录时如果用户不存在，自动创建账户并生成随机昵称和默认头像
- 订阅状态支持免费/付费会员
- 新用户昵称格式：配置前缀 + 8位随机字符（数字+字母）
- 头像上传流程：
  1. 客户端调用 GetAvatarUploadToken 获取上传凭证
  2. 客户端直接上传文件到 OSS
  3. 客户端调用 CompleteAvatarUpload 完成上传并更新数据库

## 测试

### 前提条件

确保以下服务已启动：
- MySQL (端口 3306)
- Etcd (端口 2379)
- SMS Service (端口 8001)
- OSS Service (端口 8002)
- User Service (端口 8003)

### 完整流程测试

```bash
# 1. 发送验证码
grpcurl -plaintext -proto user.proto \
  -d '{"phone": "13800138000"}' \
  127.0.0.1:8003 user.User/SendVerifyCode

# 2. 登录/注册（使用真实验证码）
grpcurl -plaintext -proto user.proto \
  -d '{"phone": "13800138000", "code": "123456"}' \
  127.0.0.1:8003 user.User/LoginOrRegister

# 3. 获取用户信息
grpcurl -plaintext -proto user.proto \
  -d '{"user_id": 123456789}' \
  127.0.0.1:8003 user.User/GetUserInfo

# 4. 更新用户信息
grpcurl -plaintext -proto user.proto \
  -d '{"user_id": 123456789, "nickname": "新昵称"}' \
  127.0.0.1:8003 user.User/UpdateUserInfo

# 5. 获取头像上传凭证
grpcurl -plaintext -proto user.proto \
  -d '{"user_id": 123456789, "file_name": "avatar.jpg"}' \
  127.0.0.1:8003 user.User/GetAvatarUploadToken

# 6. 获取订阅信息
grpcurl -plaintext -proto user.proto \
  -d '{"user_id": 123456789}' \
  127.0.0.1:8003 user.User/GetSubscription
```

### 使用测试脚本

项目提供了自动化测试脚本：

```bash
cd /Users/flf/Desktop/ProjectCode/me2/backend
./scripts/test-user.sh
```

测试脚本会自动测试：
- 发送验证码
- 登录/注册
- 获取用户信息
- 更新订阅
- 获取订阅信息
- 更新用户状态
