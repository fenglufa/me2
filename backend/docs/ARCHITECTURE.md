# 架构设计说明

## 模块化架构 (Polyrepo)

本项目采用 **Polyrepo** 模式，每个服务都是独立的 Go 模块。

### 设计原则

1. **服务独立性**: 每个服务有独立的 `go.mod`，可以单独开发、测试、部署
2. **依赖隔离**: 每个服务只包含自己需要的依赖
3. **团队协作**: 不同服务可以由不同团队负责，甚至可以是独立的代码仓库
4. **公共库复用**: `pkg` 作为独立模块，所有服务通过 Go module 引用

### 目录结构

```
backend/
├── pkg/                    # 公共库 (独立模块)
│   ├── go.mod             # module: github.com/me2/pkg
│   ├── errcode/           # 错误码
│   ├── response/          # 统一响应
│   ├── middleware/        # 中间件
│   └── utils/             # 工具函数
│
└── services/              # 微服务目录
    ├── sms/               # 短信服务 (独立模块)
    │   ├── go.mod         # module: github.com/me2/sms
    │   ├── Makefile
    │   └── rpc/
    │
    └── oss/               # OSS 服务 (独立模块)
        ├── go.mod         # module: github.com/me2/oss
        ├── Makefile
        └── rpc/
```

## 模块依赖关系

### pkg 模块

```go
// backend/pkg/go.mod
module github.com/me2/pkg

go 1.21

require (
    github.com/zeromicro/go-zero v1.6.1
)
```

### 服务模块 (以 SMS 为例)

```go
// backend/services/sms/go.mod
module github.com/me2/sms

go 1.21

require (
    github.com/me2/pkg v0.0.1  // 引用公共库
    github.com/zeromicro/go-zero v1.6.1
    // ... 其他依赖
)

// 本地开发时使用 replace
replace github.com/me2/pkg => ../../pkg
```

## Import 路径规范

### 服务内部引用

```go
// SMS 服务内部引用
import (
    "github.com/me2/sms/rpc/internal/config"
    "github.com/me2/sms/rpc/internal/svc"
    "github.com/me2/sms/rpc/sms"
)
```

### 引用公共库

```go
// 引用 pkg 公共库
import (
    "github.com/me2/pkg/errcode"
    "github.com/me2/pkg/response"
)
```

### 服务间调用

```go
// User 服务调用 SMS 服务
import (
    "github.com/me2/sms/rpc/sms"
)
```

## 开发流程

### 1. 本地开发

```bash
# 初始化所有依赖
cd backend
make init-all

# 开发某个服务
cd services/sms
make init    # 初始化依赖
make gen-rpc # 生成代码
make run     # 运行服务
```

### 2. 更新公共库 (pkg)

```bash
# 1. 修改 pkg 代码
cd backend/pkg
# ... 修改代码 ...

# 2. 各服务自动使用最新 pkg (通过 replace)
cd ../services/sms
go mod tidy  # 自动使用 ../../pkg
```

### 3. 发布公共库

当 pkg 需要发布到独立仓库时：

```bash
# 1. 将 pkg 推送到独立仓库
cd backend/pkg
git init
git remote add origin https://github.com/me2/pkg.git
git add .
git commit -m "Initial commit"
git tag v0.0.1
git push origin main --tags

# 2. 各服务更新依赖
cd ../services/sms
# 移除 replace，使用真实版本
go get github.com/me2/pkg@v0.0.1
```

### 4. 服务独立仓库

每个服务可以独立为一个仓库：

```bash
# SMS 服务独立仓库
cd backend/services/sms
git init
git remote add origin https://github.com/me2/sms-service.git
git add .
git commit -m "Initial commit"
git push origin main
```

## 优势

### ✅ 服务独立性
- 每个服务可以独立编译、测试、部署
- 服务版本独立管理
- 减少服务间耦合

### ✅ 团队协作
- 不同团队负责不同服务
- 代码权限隔离
- 减少代码冲突

### ✅ 依赖管理
- 依赖版本隔离
- 避免依赖冲突
- 构建速度更快

### ✅ 灵活部署
- 可以选择性部署某些服务
- 支持多仓库 CI/CD
- 适合微服务架构

## 注意事项

### 1. 本地开发使用 replace

```go
// 本地开发时，使用相对路径
replace github.com/me2/pkg => ../../pkg
```

### 2. 生产环境移除 replace

```go
// 生产环境使用真实版本
require github.com/me2/pkg v0.0.1
```

### 3. pkg 版本管理

- pkg 有更新时，打 tag 发布新版本
- 各服务按需更新 pkg 版本
- 保持 pkg 向后兼容

### 4. Proto 文件管理

- 每个服务的 proto 文件在服务内部
- 如果需要共享 proto，可以放在 pkg 中

## 未来扩展

当服务数量增多时，可以考虑：

1. **API Gateway**: 统一入口，路由到各个服务
2. **服务注册发现**: 使用 etcd/consul 进行服务发现
3. **配置中心**: 统一管理配置
4. **链路追踪**: 使用 Jaeger/Zipkin
5. **监控告警**: Prometheus + Grafana

## 参考资料

- [Go Modules 官方文档](https://go.dev/ref/mod)
- [go-zero 微服务实践](https://go-zero.dev/)
- [Monorepo vs Polyrepo](https://github.com/joelparkerhenderson/monorepo-vs-polyrepo)
