# goctl 代码生成规范

## 概述

本文档说明如何正确使用 goctl 工具生成符合项目命名规范的代码。

## 文件命名规范

### Go 社区标准

Go 社区推荐使用下划线分隔的文件名（snake_case），例如：
- `service_context.go`
- `send_code_logic.go`
- `verify_code_logic.go`

### goctl 的 --style 参数

goctl 支持通过 `--style` 参数控制生成文件的命名风格。

#### 可用选项

| Style 参数 | 文件命名示例 | 说明 |
|-----------|------------|------|
| `--style go_zero` 或 `--style gozero` | `service_context.go` | Go 标准风格（推荐） |
| `--style goZero` | `serviceContext.go` | 驼峰命名 |
| 默认（不带 --style） | `servicecontext.go` | 全小写驼峰命名 |

## 项目配置

### 标准 Makefile 配置

所有服务的 Makefile 都应该使用 `--style go_zero` 参数：

```makefile
# 生成 RPC 代码
.PHONY: gen-rpc
gen-rpc:
	@echo "生成 RPC 代码..."
	@cd rpc && goctl rpc protoc $(notdir $(PROTO_FILE)) --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "完成!"
```

### 命令行使用

如果需要手动生成代码，使用以下命令：

```bash
# 进入 rpc 目录
cd rpc

# 生成代码（使用 go_zero 风格）
goctl rpc protoc sms.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
```

## 生成的文件

### SMS 服务示例

使用 `--style go_zero` 生成的文件：

```
rpc/
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── logic/
│   │   ├── send_code_logic.go      ✅ 下划线命名
│   │   └── verify_code_logic.go    ✅ 下划线命名
│   ├── server/
│   │   └── sms_server.go           ✅ 下划线命名
│   └── svc/
│       └── service_context.go      ✅ 下划线命名
└── sms/
    ├── sms.pb.go
    └── sms_grpc.pb.go
```

### OSS 服务示例

使用 `--style go_zero` 生成的文件：

```
rpc/
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── logic/
│   │   ├── get_upload_url_logic.go     ✅ 下划线命名
│   │   ├── get_download_url_logic.go   ✅ 下划线命名
│   │   └── delete_file_logic.go        ✅ 下划线命名
│   ├── server/
│   │   └── oss_server.go               ✅ 下划线命名
│   └── svc/
│       └── service_context.go          ✅ 下划线命名
└── oss/
    ├── oss.pb.go
    └── oss_grpc.pb.go
```

## 常见问题

### Q1: 为什么会出现重复声明错误？

**问题**：
```
ServiceContext redeclared in this block
SendCodeLogic redeclared in this block
```

**原因**：
- 如果没有使用 `--style go_zero` 参数，goctl 会生成 `servicecontext.go`（全小写驼峰）
- 如果手动创建了 `service_context.go`（下划线分隔），就会产生两个定义相同类型的文件
- Go 编译器会报告重复声明错误

**解决方案**：
在 Makefile 中添加 `--style go_zero` 参数，确保 goctl 生成下划线命名的文件。

### Q2: 如何修复已有项目的命名问题？

**步骤**：

1. 备份现有的业务逻辑代码
2. 删除 goctl 生成的驼峰命名文件（如 `servicecontext.go`）
3. 更新 Makefile 添加 `--style go_zero` 参数
4. 重新生成代码：`make gen-rpc`
5. 恢复业务逻辑代码

### Q3: 所有服务都需要使用相同的 style 吗？

**答案**：是的。

为了保持代码风格一致性，所有服务都应该使用相同的 `--style go_zero` 参数。这样：
- 代码风格统一
- 减少认知负担
- 便于团队协作
- 符合 Go 社区标准

## 最佳实践

### 1. 统一配置

在项目根目录的 Makefile 中定义统一的代码生成规则：

```makefile
# 在根 Makefile 中
.PHONY: gen-all
gen-all:
	@cd services/sms && make gen-rpc
	@cd services/oss && make gen-rpc
	@cd services/user && make gen-rpc
```

### 2. 代码审查检查清单

在代码审查时，确保：
- [ ] 所有 Makefile 都使用 `--style go_zero` 参数
- [ ] 生成的文件名使用下划线分隔
- [ ] 没有驼峰命名的文件（如 `servicecontext.go`）
- [ ] proto 文件的 service 名称使用 PascalCase（如 `Sms`、`Oss`）

### 3. CI/CD 集成

在 CI/CD 流程中添加检查：

```bash
# 检查是否存在驼峰命名的文件
if find services -name "*context.go" | grep -v "service_context.go"; then
  echo "错误: 发现驼峰命名的文件，请使用 --style go_zero 重新生成"
  exit 1
fi
```

### 4. 新服务创建流程

创建新服务时：

1. 编写 proto 文件
2. 在 Makefile 中配置 `--style go_zero`
3. 运行 `make gen-rpc` 生成代码
4. 实现业务逻辑
5. 不要手动创建框架文件（让 goctl 生成）

## 参考资料

- [goctl 官方文档](https://go-zero.dev/docs/tutorials)
- [Go 代码规范](https://go.dev/doc/effective_go)
- 项目文档：`docs/TROUBLESHOOTING.md`
