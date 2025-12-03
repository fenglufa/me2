# 快速开始指南

本指南帮助你快速搭建 ME2 后端开发环境并运行服务。

> **快速启动（适合有经验的开发者）**:
> ```bash
> # 1. 确保 Etcd 和 Redis 运行中
> etcdctl endpoint health && redis-cli ping
>
> # 2. 启动服务（开发模式）
> cd services/sms && make run-dev  # 启动 SMS 服务
> cd services/oss && make run-dev  # 启动 OSS 服务（新终端）
>
> # 3. 验证服务注册
> ./scripts/etcd-manager.sh list
> ```

## 前置条件

- Go 1.21+
- **Etcd** (服务注册中心，必需)
- **Redis** (SMS 服务需要，必需)
- PostgreSQL (可选，用户服务需要)

## 第一步：安装依赖工具

### 安装 Go 工具链

```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装 protoc
# macOS
brew install protobuf

# Linux
apt-get install -y protobuf-compiler  # Ubuntu/Debian
yum install -y protobuf-compiler      # CentOS/RHEL

# 安装 protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 验证安装
goctl --version
protoc --version
```

### 安装 Etcd

```bash
# macOS
brew install etcd
brew services start etcd

# Linux (Docker)
docker run -d --name etcd \
  -p 2379:2379 -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  bitnami/etcd:latest

# 验证 Etcd
etcdctl version
etcdctl put test "hello"
etcdctl get test
```

### 安装 Redis (可选)

```bash
# macOS
brew install redis
brew services start redis

# Linux (Docker)
docker run -d --name redis -p 6379:6379 redis:latest

# 验证 Redis
redis-cli ping  # 应该返回 PONG
```

## 第二步：配置环境变量

```bash
# 进入项目目录
cd /Users/flf/Desktop/ProjectCode/me2/backend

# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

**最小配置**（用于开发测试）:
```bash
# Etcd 配置
ETCD_HOST=127.0.0.1:2379

# Redis 配置
REDIS_HOST=localhost:6379
REDIS_PASS=

# 阿里云配置（暂时可以填假值用于测试代码生成）
ALIYUN_ACCESS_KEY_ID=your_key
ALIYUN_ACCESS_KEY_SECRET=your_secret
ALIYUN_SMS_SIGN_NAME=测试签名
ALIYUN_SMS_TEMPLATE_CODE=SMS_123456
ALIYUN_OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
ALIYUN_OSS_BUCKET_NAME=test-bucket
```

## 第三步：初始化项目

### 方法一：一键初始化（推荐）

```bash
cd /Users/flf/Desktop/ProjectCode/me2/backend

# 一键初始化（自动生成代码 + 下载依赖）
make init-all
```

这个命令会自动：
1. 生成所有服务的 RPC 代码
2. 初始化 pkg 依赖
3. 初始化所有服务依赖

### 方法二：分步初始化

```bash
# 1. 生成所有 RPC 代码
make gen-all

# 2. 初始化所有依赖
make init-all
```

### 方法三：单独初始化某个服务

```bash
# 初始化 SMS 服务
cd services/sms
make gen-rpc  # 生成代码
make init     # 初始化依赖

# 初始化 OSS 服务
cd services/oss
make gen-rpc  # 生成代码
make init     # 初始化依赖
```

## 第四步：运行服务

### 启动 SMS 服务

```bash
cd /Users/flf/Desktop/ProjectCode/me2/backend/services/sms
make run
```

你应该看到：
```
Starting sms rpc server at 0.0.0.0:8001...
```

### 启动 OSS 服务（新终端窗口）

```bash
cd /Users/flf/Desktop/ProjectCode/me2/backend/services/oss
make run
```

你应该看到：
```
Starting oss rpc server at 0.0.0.0:8002...
```

### 验证服务注册

```bash
# 查看 SMS 服务注册
etcdctl get --prefix "sms.rpc"

# 查看 OSS 服务注册
etcdctl get --prefix "oss.rpc"
```

## 第五步：测试服务（可选）

### 使用 grpcurl 测试

```bash
# 安装 grpcurl
brew install grpcurl  # macOS
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 列出 SMS 服务的方法
grpcurl -plaintext localhost:8001 list

# 调用 SendCode 方法
grpcurl -plaintext -d '{
  "phone": "13800138000",
  "scene": "login"
}' localhost:8001 sms.Sms/SendCode
```

## 常见问题

### 1. make init 报错: cannot find module

**错误信息**:
```
cannot find module providing package github.com/me2/sms/rpc/sms
```

**原因**: proto 代码还未生成

**解决**: 先执行 `make gen-rpc` 生成代码

```bash
cd services/sms
make gen-rpc  # 先生成
make init     # 再初始化
```

### 2. Etcd 连接失败

**错误信息**:
```
failed to connect to etcd
```

**检查 Etcd 是否运行**:
```bash
# 检查进程
ps aux | grep etcd

# 检查端口
lsof -i :2379

# 重启 Etcd
brew services restart etcd  # macOS
```

### 3. Redis 连接失败

**错误信息**:
```
dial tcp [::1]:6379: connect: connection refused
```

**检查 Redis 是否运行**:
```bash
# 检查 Redis
redis-cli ping

# 启动 Redis
brew services start redis  # macOS
```

### 4. 端口被占用

**错误信息**:
```
bind: address already in use
```

**解决**:
```bash
# 查看端口占用
lsof -i :8001
lsof -i :8002

# 杀死进程
kill -9 <PID>
```

## 开发流程

### 修改 proto 文件后

```bash
# 重新生成代码
cd services/sms
make gen-rpc

# 更新依赖
make init
```

### 添加新的依赖包

```bash
# 进入服务目录
cd services/sms

# 添加依赖
go get github.com/xxx/yyy

# 更新 go.mod
make init
```

### 代码格式化和检查

```bash
# 格式化代码
make fmt-all

# 代码检查
make lint-all

# 运行测试
make test-all
```

## 下一步

- 阅读 [后端开发规范](../docs/后端开发规范.md)
- 阅读 [架构设计说明](./ARCHITECTURE.md)
- 阅读 [Etcd 使用指南](./docs/ETCD.md)
- 开始开发新服务

## 需要帮助？

- 查看服务文档: `backend/services/{service}/README.md`
- 查看 Makefile: `make help`
- 查看项目 README: `backend/README.md`
