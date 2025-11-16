.PHONY: help build run test clean install-deps fmt lint

# 项目名称
PROJECT_NAME=ldap-microservice
BINARY_NAME=$(PROJECT_NAME).exe

# 默认目标
help:
	@echo "LDAP 微服务 - Makefile 命令"
	@echo ""
	@echo "可用命令:"
	@echo "  make build          - 编译项目"
	@echo "  make run            - 编译并运行项目"
	@echo "  make test           - 运行单元测试"
	@echo "  make test-api       - 运行 API 测试"
	@echo "  make clean          - 清理编译产物"
	@echo "  make install-deps   - 安装依赖"
	@echo "  make fmt            - 格式化代码"
	@echo "  make lint           - 代码检查"
	@echo "  make help           - 显示此帮助信息"
	@echo ""

# 编译项目
build:
	@echo "编译项目..."
	go build -o $(BINARY_NAME)
	@echo "✓ 编译完成: $(BINARY_NAME)"

# 编译并运行
run: build
	@echo "启动服务..."
	./$(BINARY_NAME)

# 运行单元测试
test:
	@echo "运行单元测试..."
	go test -v ./...
	@echo "✓ 测试完成"

# 运行 API 测试
test-api:
	@echo "运行 API 测试..."
	@if exist test-api.ps1 (powershell -ExecutionPolicy Bypass -File test-api.ps1) else (echo "test-api.ps1 not found")

# 清理编译产物
clean:
	@echo "清理编译产物..."
	@if exist $(BINARY_NAME) del $(BINARY_NAME)
	@echo "✓ 清理完成"

# 安装依赖
install-deps:
	@echo "安装依赖..."
	go mod download
	go mod tidy
	@echo "✓ 依赖安装完成"

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...
	@echo "✓ 代码格式化完成"

# 代码检查
lint:
	@echo "运行代码检查..."
	@if command -v golangci-lint >nul 2>&1 (golangci-lint run) else (echo "golangci-lint not installed, skipping")
	@echo "✓ 代码检查完成"

# 完整构建流程
all: clean install-deps fmt build test
	@echo "✓ 完整构建完成"

# 开发模式（监视文件变化并重新编译）
dev:
	@echo "启动开发模式..."
	@if command -v air >nul 2>&1 (air) else (echo "air not installed, use: go install github.com/cosmtrek/air@latest")

# 显示版本信息
version:
	@echo "LDAP 微服务"
	@go version
	@echo "Go 模块:"
	@go list -m all | head -5

