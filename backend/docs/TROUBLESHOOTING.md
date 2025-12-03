# 编译问题排查指南

## 常见编译问题

### 1. grpc 版本不兼容

**错误信息**:
```
undefined: grpc.SupportPackageIsVersion9
undefined: grpc.StaticMethod
```

**原因**: proto 生成的代码需要更新版本的 grpc 包

**解决方案**:
```bash
# 更新 grpc 到最新版本
go get -u google.golang.org/grpc@latest
go mod tidy

# 重新生成 RPC 代码
make gen-rpc

# 编译
make build
```

### 2. 重复声明错误

**错误信息**:
```
ServiceContext redeclared in this block
SendCodeLogic redeclared in this block
```

**原因**: goctl 默认生成驼峰命名的文件（如 `servicecontext.go`），与我们的下划线命名文件（如 `service_context.go`）冲突

**正确解决方案**: 使用 goctl 的 `--style go_zero` 参数

goctl 支持通过 `--style` 参数控制生成的文件命名风格：
- `--style gozero` 或 `--style go_zero`: 生成下划线命名（`service_context.go`）
- 默认（不带 --style）: 生成驼峰命名（`servicecontext.go`）

现在的 Makefile 使用 `--style go_zero` 参数：
```makefile
gen-rpc:
	@echo "生成 RPC 代码..."
	@cd rpc && goctl rpc protoc $(notdir $(PROTO_FILE)) --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "完成!"
```

这样 goctl 生成的文件名就与我们的命名规范一致，不会产生重复文件。

### 3. 未找到模块

**错误信息**:
```
cannot find module providing package github.com/me2/sms/rpc/sms
```

**原因**: proto 代码还未生成

**解决方案**:
```bash
# 正确的顺序
make gen-rpc  # 1. 先生成代码
make init     # 2. 再初始化依赖
make build    # 3. 编译
```

### 4. Go 版本要求

**当前版本要求**: Go 1.24.0+

grpc v1.77.0 要求 Go 1.24+，系统会自动切换到 go1.24.10

```bash
# 检查 Go 版本
go version

# 如果版本太低，升级 Go
# macOS
brew upgrade go

# Linux
# 从官网下载最新版本
```

## 文件命名规范

### goctl 的 --style 参数

goctl 支持不同的命名风格，通过 `--style` 参数控制：

**使用 `--style go_zero` (推荐)**:
- `service_context.go` - 下划线命名的服务上下文
- `send_code_logic.go` - 下划线命名的 logic
- `sms_server.go` - 下划线命名的 server

**默认（不带 --style）**:
- `servicecontext.go` - 驼峰命名的服务上下文
- `sendcodelogic.go` - 驼峰命名的 logic
- `smsserver.go` - 驼峰命名的 server

**项目规范**: 我们统一使用 `--style go_zero` 生成下划线命名的文件，与 Go 社区标准命名规范一致。

## 日志方法

go-zero 的 Logger 支持的方法：

### 可用方法
```go
l.Debug()   // 调试日志
l.Debugf()  // 格式化调试日志
l.Info()    // 信息日志
l.Infof()   // 格式化信息日志
l.Error()   // 错误日志
l.Errorf()  // 格式化错误日志
l.Slow()    // 慢查询日志
l.Slowf()   // 格式化慢查询日志
```

### ❌ 不可用方法
```go
l.Warn()    // ❌ 不存在
l.Warnf()   // ❌ 不存在
```

**建议**: 警告信息使用 `l.Errorf()` 或 `l.Infof()`

## 完整的开发流程

```bash
# 1. 修改 proto 文件
vim rpc/sms.proto

# 2. 生成代码（自动清理重复文件）
make gen-rpc

# 3. 更新依赖
make init

# 4. 编译
make build

# 5. 运行
make run
```

## Makefile 配置

### SMS Service Makefile

```makefile
# 生成 RPC 代码
.PHONY: gen-rpc
gen-rpc:
	@echo "生成 RPC 代码..."
	@cd rpc && goctl rpc protoc $(notdir $(PROTO_FILE)) --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "完成!"
```

### OSS Service Makefile

```makefile
# 生成 RPC 代码
.PHONY: gen-rpc
gen-rpc:
	@echo "生成 RPC 代码..."
	@cd rpc && goctl rpc protoc $(notdir $(PROTO_FILE)) --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "完成!"
```

**关键点**:
- 使用 `--style go_zero` 参数确保生成下划线命名的文件
- 不需要额外的文件删除步骤
- 生成的文件直接符合项目命名规范

## 快速排查

遇到编译问题时，按以下顺序排查：

1. **检查 Go 版本** - `go version` (需要 1.24+)
2. **更新依赖** - `make init`
3. **重新生成代码** - `make gen-rpc`
4. **清理构建缓存** - `go clean -cache`
5. **重新编译** - `make build`

## 依赖版本

当前使用的主要依赖版本：

```
go 1.24.0
google.golang.org/grpc v1.77.0
google.golang.org/protobuf v1.36.10
github.com/zeromicro/go-zero v1.6.1
```

## goctl 命名风格

我们的项目使用 `--style go_zero` 参数来确保 goctl 生成的文件名符合 Go 社区的标准命名规范（下划线分隔）。

### 支持的 style 选项

goctl 支持以下命名风格：
- `--style go_zero` 或 `--style gozero`: Go 标准风格（下划线分隔）
- `--style goZero`: 驼峰命名风格
- 默认（不带 --style）: 全小写驼峰命名

### 最佳实践

1. **统一使用 go_zero 风格** - 所有服务的 Makefile 都应该使用 `--style go_zero`
2. **文档化** - 在开发规范中明确说明使用的命名风格
3. **CI/CD 检查** - 可以添加检查确保所有文件名符合命名规范
