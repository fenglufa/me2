# 服务名称
SERVICE_NAME := ai-rpc
PROTO_FILE := rpc/ai.proto

# 默认目标
.PHONY: help
help:
	@echo "可用命令:"
	@echo "  make init        - 初始化依赖"
	@echo "  make build       - 编译服务"
	@echo "  make run         - 运行服务（生产配置）"
	@echo "  make run-dev     - 运行服务（开发配置）"
	@echo "  make test        - 运行测试"
	@echo "  make gen-rpc     - 生成 RPC 代码"
	@echo "  make clean       - 清理构建文件"

# 初始化依赖
.PHONY: init
init:
	@echo "初始化依赖..."
	@go mod tidy
	@go mod download

# 编译
.PHONY: build
build:
	@echo "编译 $(SERVICE_NAME)..."
	@cd rpc && go build -o ../bin/$(SERVICE_NAME) .

# 运行
.PHONY: run
run:
	@echo "运行 $(SERVICE_NAME)..."
	@cd rpc && go run .

# 开发模式运行（使用开发配置）
.PHONY: run-dev
run-dev:
	@echo "运行 $(SERVICE_NAME) (开发模式)..."
	@cd rpc && go run . -f etc/ai-dev.yaml

# 测试
.PHONY: test
test:
	@echo "运行测试..."
	@go test -v ./...

# 生成 RPC 代码
.PHONY: gen-rpc
gen-rpc:
	@echo "生成 RPC 代码..."
	@cd rpc && goctl rpc protoc $(notdir $(PROTO_FILE)) --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "完成!"

# 清理
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -rf rpc/ai/
	@rm -rf rpc/ai_client/
