# Dialogue Service

对话服务 - 提供 AI 分身与用户的对话功能

## 功能特性

- 会话管理(创建、查询、删除)
- 流式对话(基于 gRPC Stream)
- 上下文构建(人格 + 最近事件 + 对话历史)
- 消息历史记录
- 自动生成会话标题

## 目录结构

```
dialogue/
├── Makefile              # 构建和管理命令
├── README.md            # 本文档
├── sql/                 # 数据库脚本
│   └── dialogue.sql     # 建表 SQL
└── rpc/                 # RPC 服务
    ├── dialogue.proto   # Proto 定义
    ├── etc/             # 配置文件
    │   └── dialogue-dev.yaml
    ├── internal/
    │   ├── config/      # 配置结构
    │   ├── logic/       # 业务逻辑
    │   │   ├── create_session_logic.go
    │   │   ├── get_sessions_logic.go
    │   │   ├── get_session_info_logic.go
    │   │   ├── get_messages_logic.go
    │   │   ├── delete_session_logic.go
    │   │   ├── chat_stream_logic.go
    │   │   ├── context_builder.go
    │   │   └── utils.go
    │   ├── model/       # 数据模型
    │   │   ├── dialogue_session_model.go
    │   │   └── dialogue_message_model.go
    │   ├── server/      # gRPC Server
    │   └── svc/         # Service Context
    └── dialogue_client/ # RPC Client

## 数据库表

### dialogue_sessions (对话会话表)
- id: 会话ID
- user_id: 用户ID
- avatar_id: 分身ID
- title: 会话标题
- last_message: 最后一条消息
- created_at: 创建时间
- updated_at: 更新时间

### dialogue_messages (对话消息表)
- id: 消息ID
- session_id: 会话ID
- role: 角色(user/assistant)
- content: 消息内容
- created_at: 创建时间

## 配置说明

```yaml
Name: dialogue.rpc
ListenOn: 0.0.0.0:8009

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: dialogue.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me2?charset=utf8mb4&parseTime=true&loc=Local

# 依赖的 RPC 服务
AvatarRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: avatar.rpc

EventRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: event.rpc

AiRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: ai.rpc

# 上下文配置
Context:
  MaxHistoryMessages: 10  # 最多包含的历史消息数
  MaxRecentEvents: 5      # 最多包含的最近事件数
```

## 快速开始

### 1. 初始化数据库

```bash
mysql -u root -p me2 < sql/dialogue.sql
```

### 2. 安装依赖

```bash
make init
```

### 3. 生成代码(如果修改了 proto)

```bash
make gen-rpc
```

### 4. 运行服务

开发模式:
```bash
make run-dev
```

生产模式:
```bash
make run
```

### 5. 编译

```bash
make build
```

## API 接口

### CreateSession - 创建会话
```protobuf
rpc CreateSession(CreateSessionRequest) returns (CreateSessionResponse);
```

### GetSessions - 获取会话列表
```protobuf
rpc GetSessions(GetSessionsRequest) returns (GetSessionsResponse);
```

### GetSessionInfo - 获取会话详情
```protobuf
rpc GetSessionInfo(GetSessionInfoRequest) returns (GetSessionInfoResponse);
```

### GetMessages - 获取对话历史
```protobuf
rpc GetMessages(GetMessagesRequest) returns (GetMessagesResponse);
```

### ChatStream - 流式对话
```protobuf
rpc ChatStream(ChatStreamRequest) returns (stream ChatStreamResponse);
```

### DeleteSession - 删除会话
```protobuf
rpc DeleteSession(DeleteSessionRequest) returns (DeleteSessionResponse);
```

## 核心逻辑

### 上下文构建 (ContextBuilder)

对话服务会自动构建包含以下信息的上下文:

1. **分身人格**: 从 Avatar Service 获取 6 维人格向量并转换为文字描述
2. **最近事件**: 从 Event Service 获取最近 N 条事件
3. **对话历史**: 从数据库获取当前会话的最近 N 条消息

### 流式对话流程

1. 验证会话权限
2. 构建系统 Prompt(人格 + 事件 + 历史)
3. 保存用户消息
4. 调用 AI Service 的 ChatStream 接口
5. 流式返回 AI 响应
6. 保存 AI 响应消息
7. 更新会话的最后一条消息

## 依赖服务

- **Avatar Service**: 获取分身人格信息
- **Event Service**: 获取最近事件
- **AI Service**: 调用 DeepSeek 生成对话响应
- **MySQL**: 存储会话和消息
- **Etcd**: 服务注册与发现

## 开发规范

- 所有文件名使用下划线命名(snake_case)
- 单个文件不超过 500 行,超过则拆分
- 使用 `--style go_zero` 生成代码
- 所有 RPC 方法都需要验证用户权限

## 测试

```bash
make test
```

## 代码检查

```bash
make lint
```

## 格式化代码

```bash
make fmt
```

## 清理构建文件

```bash
make clean
```

## 注意事项

1. 确保 MySQL、Etcd 已启动
2. 确保 Avatar、Event、AI 服务已启动
3. 首次运行需要执行数据库迁移
4. 配置文件中的端口不要与其他服务冲突

## 版本信息

- 版本: v1.0.0
- Go 版本: 1.24.0+
- go-zero 版本: 1.9.3+

## 作者

ME2 Team

## 许可证

MIT
