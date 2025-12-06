# Gateway Service

ME2 项目的 API 网关服务,提供统一的 HTTP API 入口,负责认证、路由转发和响应封装。

## 功能特性

- ✅ 统一 API 入口
- ✅ JWT 认证和鉴权
- ✅ 请求日志记录
- ✅ Panic 恢复
- ✅ 统一响应格式
- ✅ 路由转发到 RPC 服务
- ✅ 支持 6 个核心服务的完整 API
- ✅ 使用 pkg 公共库的中间件和工具

## 目录结构

```
gateway/
├── api/
│   ├── desc/                 # API 定义文件(按服务拆分)
│   │   ├── user.api         # 用户服务 API
│   │   ├── avatar.api       # 分身服务 API
│   │   ├── event.api        # 事件服务 API
│   │   └── diary.api        # 日记服务 API
│   ├── etc/
│   │   └── gateway.yaml     # 配置文件
│   ├── internal/
│   │   ├── config/          # 配置结构
│   │   ├── handler/         # HTTP 处理器(按服务分组)
│   │   │   ├── user/
│   │   │   ├── avatar/
│   │   │   ├── world/
│   │   │   ├── event/
│   │   │   ├── action/
│   │   │   └── diary/
│   │   ├── logic/           # 业务逻辑(按服务分组)
│   │   │   ├── user/
│   │   │   ├── avatar/
│   │   │   ├── world/
│   │   │   ├── event/
│   │   │   ├── action/
│   │   │   └── diary/
│   │   ├── middleware/      # 中间件
│   │   ├── svc/             # 服务上下文
│   │   └── types/           # 类型定义
│   ├── gateway.api          # 主 API 文件
│   └── gateway.go           # 入口文件
├── Makefile
├── go.mod
└── README.md
```

## API 分组

### 1. 用户服务 (`/api/v1/user`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/send-code` | 发送验证码 | 否 |
| POST | `/login` | 用户登录 | 否 |
| GET  | `/info` | 获取用户信息 | 是 |

### 2. 分身服务 (`/api/v1/avatar`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST  | `/create` | 创建分身 | 是 |
| GET   | `/my` | 获取我的分身 | 是 |
| GET   | `/:id` | 获取分身详情 | 是 |
| PATCH | `/:id` | 更新分身资料 | 是 |

### 3. 世界服务 (`/api/v1/world`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/maps` | 获取世界地图列表 | 是 |
| GET | `/maps/:id` | 获取世界地图详情 | 是 |
| GET | `/regions` | 获取地图的区域列表 | 是 |
| GET | `/regions/:id` | 获取区域详情 | 是 |
| GET | `/scenes` | 获取区域内的场景列表 | 是 |
| GET | `/scenes/:id` | 获取场景详情 | 是 |
| GET | `/recommend` | 根据行为类型推荐场景 | 是 |

### 4. 事件服务 (`/api/v1/events`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/timeline` | 获取事件时间线 | 是 |
| GET | `/:id` | 获取事件详情 | 是 |

### 5. 行动服务 (`/api/v1/action`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/intent` | 计算分身行动意图 | 是 |
| GET | `/history` | 获取分身行动历史 | 是 |
| GET | `/last` | 获取最近一次行动 | 是 |

### 6. 日记服务 (`/api/v1/diary`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET  | `/avatar` | 获取分身日记列表 | 是 |
| GET  | `/user` | 获取用户日记列表 | 是 |
| POST | `/user` | 创建用户日记 | 是 |
| GET  | `/:id` | 获取日记详情 | 是 |

## 配置说明

### gateway.yaml

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8888
Mode: dev

# JWT 配置
Auth:
  AccessSecret: me2-jwt-secret-key-2024
  AccessExpire: 604800  # 7天

# RPC 服务配置
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
  Timeout: 5000

AvatarRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: avatar.rpc
  Timeout: 5000

EventRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: event.rpc
  Timeout: 5000

DiaryRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: diary.rpc
  Timeout: 5000

WorldRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: world.rpc
  Timeout: 5000

ActionRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: action.rpc
  Timeout: 5000

# 日志配置
Log:
  ServiceName: gateway
  Mode: console
  Level: info
  Encoding: json
```

## 使用方式

### 1. 初始化依赖

```bash
make init
```

### 2. 生成 API 代码

```bash
make gen-api
```

### 3. 运行服务

```bash
# 开发模式
make run-dev

# 或直接运行
cd api && go run gateway.go -f etc/gateway.yaml
```

### 4. 编译

```bash
make build
```

## 开发指南

### 添加新的 API

1. **创建 API 定义文件**

在 `api/desc/` 目录下创建新的 `.api` 文件:

```api
syntax = "v1"

type (
    // 定义请求和响应类型
    ExampleRequest {
        Field string `json:"field"`
    }

    ExampleResponse {
        Data string `json:"data"`
    }
)

@server (
    group:      example
    prefix:     /api/v1/example
    middleware: Auth
)
service gateway {
    @doc "示例接口"
    @handler exampleHandler
    post /action (ExampleRequest) returns (ExampleResponse)
}
```

2. **在主 API 文件中引入**

编辑 `api/gateway.api`:

```api
import "desc/example.api"
```

3. **重新生成代码**

```bash
make gen-api
```

4. **实现 Logic 层**

在 `api/internal/logic/example/` 目录下实现业务逻辑:

```go
func (l *ExampleLogic) Example(req *types.ExampleRequest) (*types.ExampleResponse, error) {
    // 1. 调用 RPC 服务
    resp, err := l.svcCtx.ExampleRpc.Method(l.ctx, &example.Request{
        Field: req.Field,
    })
    if err != nil {
        return nil, err
    }

    // 2. 返回响应
    return &types.ExampleResponse{
        Data: resp.Data,
    }, nil
}
```

### 中间件使用

Gateway 使用 pkg 提供的中间件:

- **Auth**: JWT 认证中间件,自动验证 Token 并提取用户 ID
- **Logger**: 请求日志中间件,记录所有 HTTP 请求
- **Recovery**: Panic 恢复中间件,捕获异常并返回错误响应

在 API 定义中使用:

```api
@server (
    middleware: Auth  # 需要认证
)
service gateway {
    // ...
}
```

### 响应格式

所有 API 响应使用统一格式:

**成功响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    // 业务数据
  }
}
```

**错误响应:**
```json
{
  "code": 10001,
  "msg": "参数错误",
  "data": null
}
```

**分页响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "list": [...]
  }
}
```

## 依赖服务

Gateway 依赖以下 RPC 服务:

- **User Service** - 用户认证和信息管理
- **Avatar Service** - 分身管理
- **World Service** - 世界地图和场景管理
- **Event Service** - 事件管理
- **Action Service** - 行动逻辑管理
- **Diary Service** - 日记管理

确保这些服务已启动并注册到 Etcd。

## 开发规范

1. **API 定义按服务拆分** - 每个服务一个 `.api` 文件
2. **Handler 和 Logic 按服务分组** - 使用目录区分不同服务
3. **使用 pkg 公共库** - 错误码、响应格式、中间件等
4. **代码文件使用下划线命名** - 如 `send_code_logic.go`
5. **单文件不超过 500 行** - 超过则拆分

## 常见问题

### 1. RPC 服务连接失败

检查:
- Etcd 是否启动
- RPC 服务是否已注册到 Etcd
- 配置文件中的 Etcd 地址是否正确

### 2. JWT 认证失败

检查:
- Token 是否正确传递 (Header: `Authorization: Bearer <token>`)
- Token 是否过期
- AccessSecret 配置是否一致

### 3. 跨域问题

在 gateway.go 中添加 CORS 中间件:

```go
server.Use(rest.WithCors())
```

## 性能优化

1. **连接池配置** - RPC 客户端使用连接池
2. **超时设置** - 合理设置 RPC 调用超时时间
3. **日志级别** - 生产环境使用 `info` 或 `warn` 级别
4. **限流** - 使用 go-zero 内置限流中间件

## 监控

Gateway 提供以下监控指标:

- HTTP 请求数
- 请求响应时间
- 错误率
- RPC 调用统计

使用 Prometheus + Grafana 进行监控。

## 版本

- **v1.0.0** - 初始版本,实现基础 API 网关功能

## 作者

ME2 Team

## 许可

MIT License
