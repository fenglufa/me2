# Avatar Service

分身服务，提供分身创建、信息管理、人格系统和头像上传功能。

## 目录

- [功能](#功能)
- [技术栈](#技术栈)
- [配置](#配置)
- [数据库](#数据库)
- [使用](#使用)
- [API](#api)
  - [CreateAvatar - 创建分身](#createavatar---创建分身)
  - [GetMyAvatar - 获取我的分身](#getmyavatar---获取我的分身)
  - [GetAvatarInfo - 获取分身详情](#getavatarinfo---获取分身详情)
  - [UpdateAvatarProfile - 更新分身资料](#updateavatarprofile---更新分身资料)
  - [GetAvatarUploadToken - 获取头像上传凭证](#getavataruploadtoken---获取头像上传凭证)
  - [CompleteAvatarUpload - 完成头像上传](#completeavatarupload---完成头像上传)
- [人格生成系统](#人格生成系统)
- [头像上传流程](#头像上传流程)
- [依赖服务](#依赖服务)
- [注意事项](#注意事项)

## 功能

- 创建分身 (CreateAvatar) - 收集用户信息并自动生成 6 维人格
- 获取我的分身 (GetMyAvatar) - 查询用户是否已创建分身
- 获取分身详情 (GetAvatarInfo) - 获取完整的分身信息
- 更新分身资料 (UpdateAvatarProfile) - 更新昵称、头像
- 获取头像上传凭证 (GetAvatarUploadToken) - 调用 OSS 服务获取上传凭证
- 完成头像上传 (CompleteAvatarUpload) - 完成上传并更新分身头像

## 技术栈

- go-zero RPC
- MySQL (分身数据存储)
- Etcd (服务注册与发现)
- 雪花算法 (生成 10 位分身 ID)
- 基于人口统计学的人格生成算法

## 配置

配置文件: `rpc/etc/avatar-dev.yaml`

```yaml
Name: avatar.rpc
ListenOn: 0.0.0.0:8004

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: avatar.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me2?charset=utf8mb4&parseTime=true&loc=Local

OssRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: oss.rpc

MachineID: 1  # 雪花算法机器ID (0-1023)
```

## 数据库

### 创建数据库表

```bash
mysql -u root -p < rpc/avatar.sql
```

### 表结构

```sql
CREATE TABLE avatars (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    avatar_id BIGINT UNIQUE NOT NULL COMMENT '分身ID（10位数字）',
    user_id BIGINT UNIQUE NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar_url VARCHAR(500) DEFAULT '',
    gender TINYINT NOT NULL,
    birth_date DATE NOT NULL,
    occupation VARCHAR(50) DEFAULT '',
    marital_status TINYINT DEFAULT 1,

    -- 6维人格
    warmth TINYINT DEFAULT 50,
    adventurous TINYINT DEFAULT 50,
    social TINYINT DEFAULT 50,
    creative TINYINT DEFAULT 50,
    calm TINYINT DEFAULT 50,
    energetic TINYINT DEFAULT 50,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_user_id (user_id),
    INDEX idx_avatar_id (avatar_id)
);
```

## 使用

```bash
# 1. 生成 RPC 代码 (必须先执行)
make gen-rpc

# 2. 初始化依赖
make init

# 3. 创建数据库表
mysql -u root -p < rpc/avatar.sql

# 4. 编译
make build

# 5. 运行（开发模式）
make run-dev

# 6. 测试
make test
```

**注意**: 必须先执行 `make gen-rpc` 生成 proto 代码，再执行 `make init` 初始化依赖。

## API

### CreateAvatar - 创建分身

请求:
```protobuf
message CreateAvatarRequest {
  int64 user_id = 1;           // 用户ID
  string nickname = 2;         // 昵称
  string avatar_url = 3;       // 头像URL
  int32 gender = 4;            // 性别 1:男 2:女 3:其他
  string birth_date = 5;       // 出生日期 YYYY-MM-DD
  string occupation = 6;       // 职业
  int32 marital_status = 7;    // 婚姻状态 1:单身 2:恋爱中 3:已婚 4:其他
}
```

响应:
```protobuf
message CreateAvatarResponse {
  int64 avatar_id = 1;
  PersonalityInfo personality = 2;  // 生成的人格信息
}
```

### GetMyAvatar - 获取我的分身

请求:
```protobuf
message GetMyAvatarRequest {
  int64 user_id = 1;
}
```

响应:
```protobuf
message GetMyAvatarResponse {
  bool has_avatar = 1;         // 是否有分身
  AvatarInfo avatar = 2;       // 分身信息（如果有）
}
```

### GetAvatarInfo - 获取分身详情

请求:
```protobuf
message GetAvatarInfoRequest {
  int64 avatar_id = 1;
}
```

响应:
```protobuf
message GetAvatarInfoResponse {
  AvatarInfo avatar = 1;
}
```

### UpdateAvatarProfile - 更新分身资料

请求:
```protobuf
message UpdateAvatarProfileRequest {
  int64 avatar_id = 1;
  string nickname = 2;         // 昵称（可选）
  string avatar_url = 3;       // 头像URL（可选）
}
```

响应:
```protobuf
message UpdateAvatarProfileResponse {
  bool success = 1;
}
```

### GetAvatarUploadToken - 获取头像上传凭证

请求:
```protobuf
message GetAvatarUploadTokenRequest {
  int64 avatar_id = 1;         // 分身ID
  string file_name = 2;        // 文件名
}
```

响应:
```protobuf
message GetAvatarUploadTokenResponse {
  string host = 1;             // OSS 上传地址
  string access_key_id = 2;    // AccessKeyId
  string policy = 3;           // 上传策略（Base64）
  string signature = 4;        // 签名
  string key = 5;              // 文件路径
  int64 expire = 6;            // 过期时间戳
  string domain = 7;           // 文件访问域名
  string complete_token = 8;   // 完成上传的验证 Token
}
```

### CompleteAvatarUpload - 完成头像上传

请求:
```protobuf
message CompleteAvatarUploadRequest {
  int64 avatar_id = 1;         // 分身ID
  string key = 2;              // 文件 Key
  string complete_token = 3;   // 验证 Token
}
```

响应:
```protobuf
message CompleteAvatarUploadResponse {
  bool success = 1;
  string avatar_url = 2;       // 头像访问 URL
}
```

## 人格生成系统

创建分身时，系统会基于用户提供的人口统计学信息（年龄、职业、婚姻状态、性别）自动生成 6 维人格：

### 6 维人格

1. **情绪温度 (Warmth)**: 0-100，温暖 ↔ 冷静
2. **冒险倾向 (Adventurous)**: 0-100，冒险 ↔ 安全
3. **人际能量 (Social)**: 0-100，社交 ↔ 独处
4. **创造性 (Creative)**: 0-100，创造 ↔ 结构
5. **情绪稳定性 (Calm)**: 0-100，沉稳 ↔ 敏感
6. **生活动力 (Energetic)**: 0-100，活力 ↔ 温和

### 生成规则

所有维度从 50（中性）开始，然后根据以下因素调整：

**年龄影响**:
- < 25 岁: Adventurous +15, Energetic +10, Social +10
- 25-35 岁: Adventurous +5, Energetic +5
- 35-50 岁: Calm +10, Warmth +5
- > 50 岁: Calm +15, Warmth +10, Adventurous -10

**职业影响**:
- 创意类: Creative +20, Adventurous +10
- 技术类: Creative +10, Calm +10
- 商务类: Social +15, Energetic +10
- 教育类: Warmth +15, Calm +10
- 医疗类: Warmth +15, Calm +15

**婚姻状态影响**:
- 单身: Adventurous +5, Social +5
- 恋爱中: Warmth +10, Social +5
- 已婚: Warmth +10, Calm +10

**性别影响**（轻微）:
- 女性: Warmth +5
- 男性: Adventurous +5

所有值最终会被规范化到 0-100 范围内。

## 头像上传流程

分身可以通过调用头像上传接口更新头像。整个流程采用**两阶段上传模式**：

### 1. 获取上传凭证

客户端调用 `GetAvatarUploadToken` 获取上传所需的凭证：

```bash
grpcurl -plaintext -proto avatar.proto \
  -d '{"avatar_id": 123, "file_name": "avatar.jpg"}' \
  127.0.0.1:8004 avatar.Avatar/GetAvatarUploadToken
```

### 2. 客户端直接上传到 OSS

客户端使用获取的凭证，直接通过 HTTP POST 上传文件到 OSS。

### 3. 完成上传并更新数据库

客户端上传成功后，调用 `CompleteAvatarUpload` 通知服务端：

```bash
grpcurl -plaintext -proto avatar.proto \
  -d '{
    "avatar_id": 123,
    "key": "avatar/filename.jpg",
    "complete_token": "jwt-token"
  }' \
  127.0.0.1:8004 avatar.Avatar/CompleteAvatarUpload
```

### 流程说明

1. **安全性**：Avatar 服务调用 OSS 服务生成临时上传凭证
2. **性能**：客户端直接上传到 OSS，不经过应用服务器
3. **验证**：完成上传时通过 `complete_token` 验证
4. **原子性**：只有在 OSS 验证成功后，才会更新数据库

## 依赖服务

- **OSS Service**: 头像上传凭证和完成上传
- **MySQL**: 分身数据存储
- **Etcd**: 服务注册与发现

## 注意事项

- 每个用户只能创建一个分身（user_id 唯一约束）
- 分身 ID (avatar_id) 使用雪花算法生成 10 位数字，与用户 ID (9 位) 区分
- 机器 ID (MachineID) 需要在配置文件中设置，范围 0-1023，确保分布式环境下 ID 唯一性
- 分身创建时会自动生成基于人口统计学的 6 维人格
- 人格值范围为 0-100，所有维度从 50 开始调整
- 头像上传流程：
  1. 客户端调用 GetAvatarUploadToken 获取上传凭证
  2. 客户端直接上传文件到 OSS
  3. 客户端调用 CompleteAvatarUpload 完成上传并更新数据库
- 当前 MVP 版本不包含成长系统和关系系统
