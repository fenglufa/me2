# Diary Service (日记服务)

日记服务负责管理分身日记和用户日记的双向互动。

## 功能特性

### 1. 分身日记
- 每日自动生成分身日记（由 Scheduler 调度）
- 基于当天事件总结生成个性化日记
- 调用 AI Service 生成日记内容
- 记录分身的心情、标签等信息

### 2. 用户日记
- 用户写日记
- AI 自动进行情绪分析（-100 到 100）
- 分身自动生成个性化回应
- 情绪分析结果可用于更新关系分数

### 3. 日记管理
- 按日期查询日记列表
- 分页查询支持
- 日期范围筛选
- 日记统计（总数、字数、日期范围等）

###### 请注意，当前日记生成并没有使用事件内容，只是使用了标题，因为内容较多，会造成超时和 AI的 token 消耗过大，应该想办法解决

## 服务架构

```
diary/
├── rpc/                          # RPC 服务
│   ├── diary.proto              # protobuf 定义
│   ├── internal/
│   │   ├── config/              # 配置
│   │   ├── logic/               # 业务逻辑
│   │   │   ├── generate_avatar_diary_logic.go
│   │   │   ├── create_user_diary_logic.go
│   │   │   ├── get_avatar_diary_list_logic_impl.go
│   │   │   ├── get_user_diary_list_logic_impl.go
│   │   │   ├── get_diary_stats_logic_impl.go
│   │   │   └── diary_helper.go
│   │   ├── model/               # 数据模型
│   │   │   ├── diaries_model.go
│   │   │   ├── diaries_model_ext.go
│   │   │   └── diaries.sql
│   │   ├── server/              # gRPC 服务器
│   │   └── svc/                 # 服务上下文
│   └── diary.go
├── Makefile
└── README.md
```

## 数据库表

### diaries 表
```sql
CREATE TABLE `diaries` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `avatar_id` BIGINT NOT NULL,
    `type` ENUM('avatar', 'user') NOT NULL,
    `date` DATE NOT NULL,
    `title` VARCHAR(200) DEFAULT '',
    `content` TEXT NOT NULL,
    `mood` VARCHAR(50) DEFAULT '',
    `tags` VARCHAR(500) DEFAULT '',
    `reply_content` TEXT,
    `emotion_score` INT DEFAULT 0,
    `is_important` TINYINT DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_avatar_date` (`avatar_id`, `date`),
    INDEX `idx_avatar_type` (`avatar_id`, `type`),
    UNIQUE KEY `uk_avatar_type_date` (`avatar_id`, `type`, `date`)
);
```

## API 接口

### 1. GenerateAvatarDiary
生成分身日记（由 Scheduler 调用）

**请求:**
```protobuf
message GenerateAvatarDiaryRequest {
  int64 avatar_id = 1;
  string date = 2;  // 可选，默认今天
}
```

**响应:**
```protobuf
message GenerateAvatarDiaryResponse {
  int64 diary_id = 1;
  string title = 2;
  string content = 3;
  string mood = 4;
  repeated string tags = 5;
}
```

### 2. CreateUserDiary
创建用户日记

**请求:**
```protobuf
message CreateUserDiaryRequest {
  int64 avatar_id = 1;
  string title = 2;
  string content = 3;
  repeated string tags = 4;
  bool is_important = 5;
}
```

**响应:**
```protobuf
message CreateUserDiaryResponse {
  int64 diary_id = 1;
  string reply_content = 2;  // 分身回应
  int32 emotion_score = 3;   // 情绪分数
}
```

### 3. GetAvatarDiaryList
获取分身日记列表

**请求:**
```protobuf
message GetAvatarDiaryListRequest {
  int64 avatar_id = 1;
  int32 page = 2;
  int32 page_size = 3;
  string start_date = 4;  // 可选
  string end_date = 5;    // 可选
}
```

### 4. GetUserDiaryList
获取用户日记列表（同上）

### 5. GetDiaryStats
获取日记统计

**响应:**
```protobuf
message GetDiaryStatsResponse {
  int64 total_diaries = 1;
  int64 avatar_diaries = 2;
  int64 user_diaries = 3;
  int32 consecutive_days = 4;
  int64 total_words = 5;
  string first_diary_date = 6;
  string last_diary_date = 7;
}
```

## 工作流程

### 分身日记生成流程
```
Scheduler 定时触发 (每晚 22:00)
    ↓
Diary Service.GenerateAvatarDiary
    ↓
查询今天的事件 (Event Service)
    ↓
调用 AI Service 生成日记
    ↓
保存日记到数据库
    ↓
返回日记内容
```

### 用户日记创建流程
```
用户提交日记
    ↓
Diary Service.CreateUserDiary
    ↓
保存日记到数据库
    ↓
调用 AI Service 进行情绪分析
    ↓
调用 AI Service 生成分身回应
    ↓
更新日记（添加回应和情绪分数）
    ↓
返回结果
```

## 依赖服务

- **Event Service**: 查询分身当天的事件列表
- **AI Service**: 生成日记内容、情绪分析、生成回应
- **Avatar Service**: 获取分身信息（未来可能用于更新关系分数）

## 配置说明

### diary-dev.yaml
```yaml
Name: diary.rpc
ListenOn: 0.0.0.0:8011

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: diary.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me?charset=utf8mb4&parseTime=true&loc=Local

EventRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: event.rpc

AIRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: ai.rpc

AvatarRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: avatar.rpc
```

## 使用方法

### 编译
```bash
make build
```

### 运行（开发模式）
```bash
make run-dev
```

### 生成 RPC 代码
```bash
make gen-rpc
```

### 测试
```bash
make test
```

## AI Prompt 模板

需要在 AI Service 中添加以下 prompt 模板：

### 1. diary_generation
用于生成分身日记

**变量:**
- `date`: 日期
- `event_summary`: 事件摘要

**输出格式:**
```
标题：[日记标题]
心情：[心情]
标签：[标签1,标签2,标签3]
内容：[日记正文]
```

### 2. emotion_analysis
用于分析用户日记的情绪

**变量:**
- `diary_content`: 日记内容

**输出格式:**
```
[情绪分数: -100 到 100 的整数]
```

### 3. diary_reply
用于生成分身对用户日记的回应

**变量:**
- `diary_title`: 日记标题
- `diary_content`: 日记内容
- `emotion_score`: 情绪分数

**输出格式:**
```
[分身的回应文本]
```

## 注意事项

1. 每个分身每天只能有一篇分身日记（通过唯一索引保证）
2. 用户日记可以每天写多篇
3. 情绪分数范围：-100（非常消极）到 100（非常积极）
4. 日记生成依赖当天有事件，如果没有事件会跳过生成
5. AI 调用失败不会影响日记保存，但会缺少相应的内容

## 未来优化

- [ ] 实现连续记录天数计算
- [ ] 集成 Memory Service 向量化重要日记
- [ ] 集成 Notify Service 推送日记通知
- [ ] 实现日记搜索功能
- [ ] 实现日记标签系统
- [ ] 根据情绪分数更新 Avatar 关系分数
