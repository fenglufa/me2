.PHONY: help init gen-api build run run-dev test clean

help:
	@echo "可用命令:"
	@echo "  make init      - 初始化依赖"
	@echo "  make gen-api   - 生成 API 代码"
	@echo "  make build     - 编译服务"
	@echo "  make run       - 运行服务"
	@echo "  make run-dev   - 运行服务(开发模式)"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"

init:
	@echo "初始化依赖..."
	@go mod tidy
	@go mod download

gen-api:
	@echo "生成 API 代码..."
	@cd api && goctl api go --api gateway.api --dir . --style go_zero

build:
	@echo "编译 gateway..."
	@cd api && go build -o ../bin/gateway gateway.go

run:
	@echo "运行 gateway..."
	@cd api && go run gateway.go -f etc/gateway.yaml

run-dev:
	@echo "运行 gateway (开发模式)..."
	@cd api && go run gateway.go -f etc/gateway.yaml

test:
	@echo "运行测试..."
	@go test -v ./...

clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
