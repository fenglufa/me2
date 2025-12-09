# Avatar Service

分身服务，提供分身创建、信息管理和头像上传功能。

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
- [头像上传流程](#头像上传流程)
- [依赖服务](#依赖服务)
- [注意事项](#注意事项)

## 功能

- 创建分身 (CreateAvatar) - 收集用户信息并创建分身
- 获取我的分身 (GetMyAvatar) - 查询用户是否已创建分身
- 获取分身详情 (GetAvatarInfo) - 获取完整的分身信息
- 更新分身资料 (UpdateAvatarProfile) - 更新昵称、头像等基本信息
- 获取头像上传凭证 (GetAvatarUploadToken) - 调用 OSS 服务获取上传凭证
- 完成头像上传 (CompleteAvatarUpload) - 完成上传并更新分身头像

## 技术栈

- go-zero RPC
- MySQL (分身数据存储)
- Etcd (服务注册与发现)
- 雪花算法 (生成 10 位分身 ID)

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

## 分身生命周期

### 创建分身时的自动化流程

当用户创建分身后，Avatar 服务会自动触发以下操作（均为异步调用，失败不影响分身创建）：

#### 1. 启用自动调度
调用 **Scheduler Service** 的 `EnableAvatarSchedule` 接口，为分身启用自动调度：
- 调度间隔：2-6 小时随机
- 调度器每 60 秒自动扫描数据库，到时间后自动触发分身行动
- 分身将持续自主"生活"，无需外部触发

#### 2. 触发首次行动
调用 **Action Service** 的 `ScheduleAction` 接口，立即为新创建的分身触发第一次行动：
- 目的：让用户第一时间感受到分身在新世界中的活动
- 行为示例："刚进入新的世界，开始熟悉环境" 或 "开始探索周围的场景"
- 后续行动将由 Scheduler Service 自动调度

#### 3. 代码实现

```go
// 自动启用调度（异步调用，失败不影响分身创建）
go func() {
    scheduleCtx := context.Background()
    _, scheduleErr := l.svcCtx.SchedulerRpc.EnableAvatarSchedule(scheduleCtx, &scheduler.EnableAvatarScheduleRequest{
        AvatarId: avatarId,
    })
    if scheduleErr != nil {
        l.Errorf("自动启用分身调度失败 (avatar_id=%d): %v", avatarId, scheduleErr)
    } else {
        l.Infof("已为分身 %d 自动启用调度", avatarId)
    }
}()

// 立即触发首次行动（异步调用，失败不影响分身创建）
go func() {
    actionCtx := context.Background()
    _, actionErr := l.svcCtx.ActionRpc.ScheduleAction(actionCtx, &action_client.ScheduleActionRequest{
        AvatarId: avatarId,
    })
    if actionErr != nil {
        l.Errorf("触发分身首次行动失败 (avatar_id=%d): %v", avatarId, actionErr)
    } else {
        l.Infof("已为分身 %d 触发首次行动", avatarId)
    }
}()
```

### 为什么使用双异步调用？

1. **不阻塞用户体验**：分身创建成功后立即返回，不等待调度和行动完成
2. **容错性**：即使调度或行动服务暂时不可用，也不影响分身创建
3. **职责分离**：
   - **Scheduler Service** 负责"何时调度"（2-6 小时自动触发）
   - **Action Service** 负责"做什么"（决策行为逻辑）
   - **Avatar Service** 只负责分身数据管理
4. **首次体验优化**：用户不必等待 2-6 小时才看到分身的第一个动作

## 依赖服务

- **OSS Service**: 头像上传凭证和完成上传
- **Scheduler Service**: 自动调度分身行为（创建分身时自动启用）
- **Action Service**: 行为决策和执行（创建分身时立即触发首次行动）
- **MySQL**: 分身数据存储
- **Etcd**: 服务注册与发现

## 注意事项

- 每个用户只能创建一个分身（user_id 唯一约束）
- 分身 ID (avatar_id) 使用雪花算法生成 10 位数字，与用户 ID (9 位) 区分
- 机器 ID (MachineID) 需要在配置文件中设置，范围 0-1023，确保分布式环境下 ID 唯一性
- 头像上传流程：
  1. 客户端调用 GetAvatarUploadToken 获取上传凭证
  2. 客户端直接上传文件到 OSS
  3. 客户端调用 CompleteAvatarUpload 完成上传并更新数据库
