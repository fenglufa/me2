# 项目名称
PROJECT_NAME := me2

# 服务列表
SERVICES := sms oss

.PHONY: help
help:
	@echo "可用命令:"
	@echo "  make init-all    - 初始化所有服务依赖"
	@echo "  make build-all   - 编译所有服务"
	@echo "  make test-all    - 测试所有服务"
	@echo "  make fmt-all     - 格式化所有代码"
	@echo "  make lint-all    - 检查所有代码"
	@echo "  make clean-all   - 清理所有构建文件"
	@echo "  make gen-all     - 生成所有 RPC 代码"

# 初始化所有服务依赖
.PHONY: init-all
init-all: gen-all
	@echo "初始化 pkg 依赖..."
	@cd pkg && go mod tidy && go mod download
	@for service in $(SERVICES); do \
		echo "初始化 $$service 依赖..."; \
		cd services/$$service && make init && cd ../..; \
	done

# 编译所有服务
.PHONY: build-all
build-all:
	@for service in $(SERVICES); do \
		echo "编译 $$service..."; \
		cd services/$$service && make build && cd ../..; \
	done

# 测试所有服务
.PHONY: test-all
test-all:
	@for service in $(SERVICES); do \
		echo "测试 $$service..."; \
		cd services/$$service && make test && cd ../..; \
	done

# 格式化所有代码
.PHONY: fmt-all
fmt-all:
	@echo "格式化所有代码..."
	@gofmt -w .

# 检查所有代码
.PHONY: lint-all
lint-all:
	@echo "检查所有代码..."
	@golint ./...
	@go vet ./...

# 清理所有构建文件
.PHONY: clean-all
clean-all:
	@for service in $(SERVICES); do \
		echo "清理 $$service..."; \
		cd services/$$service && make clean && cd ../..; \
	done

# 生成所有 RPC 代码
.PHONY: gen-all
gen-all:
	@for service in $(SERVICES); do \
		echo "生成 $$service RPC 代码..."; \
		cd services/$$service && make gen-rpc && cd ../..; \
	done
