# ME2 Backend Services

AI 分身宇宙后端服务

## 项目结构

```
backend/
├── services/           # 微服务目录 (每个服务独立 go.mod)
│   ├── sms/           # 短信服务
│   │   ├── go.mod     # 独立模块
│   │   ├── Makefile
│   │   └── rpc/
│   └── oss/           # 对象存储服务
│       ├── go.mod     # 独立模块
│       ├── Makefile
│       └── rpc/
├── pkg/               # 公共库 (独立模块)
│   ├── go.mod         # 独立模块
│   └── README.md
├── Makefile           # 项目管理文件
└── .env.example       # 环境变量模板
```

## 模块化设计

本项目采用 **Polyrepo** 模式，每个服务都是独立的 Go 模块：

- ✅ 服务完全独立，可单独开发、测试、部署
- ✅ 依赖隔离，每个服务只包含需要的包
- ✅ 适合多人协作，每个服务可以是独立仓库
- ✅ pkg 作为独立模块，统一管理公共代码

## 已完成服务

### 1. SMS Service (短信服务)
- 端口: 8001
- 功能: 发送验证码、验证验证码
- 依赖: 阿里云短信服务、Redis、Etcd

### 2. OSS Service (对象存储服务)
- 端口: 8002
- 功能: 获取上传 URL、获取下载 URL、删除文件
- 依赖: 阿里云 OSS、Etcd

## 快速开始

详细的环境搭建和运行指南请查看 [快速开始指南](./QUICKSTART.md)

### 1. 环境准备

```bash
# 安装 Go 1.21+
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装 protoc
# macOS
brew install protobuf

# 安装 protoc-gen-go 和 protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 安装 Etcd (服务注册中心)
# macOS
brew install etcd
brew services start etcd

# Linux/Docker 安装方式见 docs/ETCD.md
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入真实配置
vim .env
```

### 3. 初始化项目

```bash
# 一键初始化 (会自动生成代码 + 初始化依赖)
make init-all

# 或者分步执行
make gen-all   # 1. 生成所有 RPC 代码
make init-all  # 2. 初始化所有服务依赖

# 或者单独初始化某个服务
cd services/sms
make gen-rpc  # 先生成代码
make init     # 再初始化依赖
```

**重要**: 必须先执行 `make gen-rpc` 生成 proto 代码，再执行 `make init` 初始化依赖。
**提示**: 执行 `make init-all` 会自动先执行 `make gen-all`，无需手动生成代码。

### 4. 运行服务

```bash
# 运行 SMS 服务
cd services/sms && make run

# 运行 OSS 服务
cd services/oss && make run
```

## 开发规范

请参考项目根目录的 [后端开发规范](../docs/后端开发规范.md)

## 相关文档

- [快速开始指南](./QUICKSTART.md) - 新手入门必读
- [架构设计说明](./ARCHITECTURE.md) - Polyrepo 模块化架构
- [Etcd 使用指南](./docs/ETCD.md) - 服务注册与发现
- [编译问题排查](./docs/TROUBLESHOOTING.md) - 常见编译问题解决方案
- [pkg 公共库](./pkg/README.md) - 公共库使用说明

## 常用命令

```bash
# 编译所有服务
make build-all

# 测试所有服务
make test-all

# 格式化代码
make fmt-all

# 代码检查
make lint-all

# 清理构建文件
make clean-all
```

## 服务端口分配

| 服务 | 端口 | 状态 |
|------|------|------|
| SMS Service | 8001 | ✅ 已完成 |
| OSS Service | 8002 | ✅ 已完成 |
| User Service | 8003 | 待开发 |
| Avatar Service | 8004 | 待开发 |
| AI Service | 8005 | 待开发 |
| Event Service | 8006 | 待开发 |
| Action Service | 8007 | 待开发 |
| World Service | 8008 | 待开发 |
| Dialogue Service | 8009 | 待开发 |
| Diary Service | 8010 | 待开发 |
| Memory Service | 8011 | 待开发 |
| ImageGen Service | 8012 | 待开发 |
| Notify Service | 8013 | 待开发 |

## 依赖服务

- **Etcd** - 服务注册与发现
- **PostgreSQL 14+** - 关系型数据库
- **Redis 7+** - 缓存和临时数据存储
- **阿里云短信服务** - 验证码发送
- **阿里云 OSS** - 对象存储
- **Deepseek API** - AI 服务

## 下一步

- [ ] 开发 User Service (用户服务)
- [ ] 开发 Avatar Service (分身服务)
- [ ] 开发 AI Service (AI 统一服务)
- [ ] 搭建数据库
- [ ] 配置 Docker Compose
