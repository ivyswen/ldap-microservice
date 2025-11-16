# Windows 环境测试指南

本指南介绍如何在 Windows 环境下测试 LDAP 微服务项目。

## 目录

1. [环境变量配置](#环境变量配置)
2. [使用 .env 文件](#使用-env-文件)
3. [启动服务](#启动服务)
4. [测试 API 端点](#测试-api-端点)
5. [常见问题](#常见问题)

## 环境变量配置

### 方法 1: 使用 PowerShell 设置环境变量

在 PowerShell 中临时设置环境变量（仅在当前会话有效）：

```powershell
# 设置 LDAP 连接参数
$env:LDAP_URL = "ldap://ldap.example.com:389"
$env:LDAP_BIND_DN = "cn=admin,dc=example,dc=com"
$env:LDAP_BIND_PASSWORD = "password123"
$env:LDAP_USER_BASE = "dc=example,dc=com"
$env:LDAP_USER_FILTER = "(uid=%s)"
$env:SERVICE_PORT = "8080"

# 验证环境变量是否设置成功
$env:LDAP_URL
```

### 方法 2: 永久设置环境变量

使用 Windows 系统设置永久设置环境变量：

1. 打开 **系统属性** → **环境变量**
2. 点击 **新建** 添加以下变量：
   - `LDAP_URL`: `ldap://ldap.example.com:389`
   - `LDAP_BIND_DN`: `cn=admin,dc=example,dc=com`
   - `LDAP_BIND_PASSWORD`: `password123`
   - `LDAP_USER_BASE`: `dc=example,dc=com`
   - `LDAP_USER_FILTER`: `(uid=%s)`
   - `SERVICE_PORT`: `8080`

3. 点击 **确定** 保存
4. 重启 PowerShell 或命令提示符使变量生效

### 方法 3: 使用批处理脚本

创建 `set-env.bat` 文件：

```batch
@echo off
REM 设置 LDAP 连接参数
set LDAP_URL=ldap://ldap.example.com:389
set LDAP_BIND_DN=cn=admin,dc=example,dc=com
set LDAP_BIND_PASSWORD=password123
set LDAP_USER_BASE=dc=example,dc=com
set LDAP_USER_FILTER=(uid=%s)
set SERVICE_PORT=8080

echo Environment variables set successfully
```

然后在命令提示符中运行：

```batch
set-env.bat
```

## 使用 .env 文件

### 推荐方法：使用 .env 文件

这是最方便的方法，特别是在开发环境中。

#### 步骤 1: 创建 .env 文件

复制 `.env.example` 文件为 `.env`：

```powershell
Copy-Item .env.example .env
```

#### 步骤 2: 编辑 .env 文件

使用任何文本编辑器打开 `.env` 文件，修改配置参数：

```env
SERVICE_PORT=8080
LDAP_URL=ldap://ldap.example.com:389
LDAP_BIND_DN=cn=admin,dc=example,dc=com
LDAP_BIND_PASSWORD=password123
LDAP_USER_BASE=dc=example,dc=com
LDAP_USER_FILTER=(uid=%s)
LDAP_USE_LDAPS=0
LDAP_USE_STARTTLS=0
LDAP_INSECURE_SKIP_VERIFY=0
```

#### 步骤 3: 启动服务

程序启动时会自动从 `.env` 文件加载配置：

```powershell
.\ldap-microservice.exe
```

### .env 文件优势

- ✅ 无需手动设置环境变量
- ✅ 配置集中管理
- ✅ 易于在不同环境间切换
- ✅ 不会污染系统环境变量
- ✅ 支持注释和文档

## 启动服务

### 编译项目

```powershell
# 进入项目目录
cd d:\GolandProjects\ldap-microservice

# 编译项目
go build -o ldap-microservice.exe

# 或使用 make（如果安装了 make）
make build
```

### 运行服务

#### 方法 1: 直接运行

```powershell
.\ldap-microservice.exe
```

#### 方法 2: 使用 PowerShell 脚本

创建 `run.ps1` 脚本：

```powershell
# 设置环境变量
$env:LDAP_URL = "ldap://ldap.example.com:389"
$env:LDAP_BIND_DN = "cn=admin,dc=example,dc=com"
$env:LDAP_BIND_PASSWORD = "password123"
$env:LDAP_USER_BASE = "dc=example,dc=com"
$env:LDAP_USER_FILTER = "(uid=%s)"
$env:SERVICE_PORT = "8080"

# 启动服务
.\ldap-microservice.exe
```

运行脚本：

```powershell
.\run.ps1
```

#### 方法 3: 使用 go run

```powershell
go run .
```

## 测试 API 端点

### 前置条件

- 服务已启动（默认监听 `http://localhost:8080`）
- 有效的 LDAP 服务器连接

### 测试工具

#### 工具 1: PowerShell Invoke-RestMethod

这是 Windows 原生工具，无需安装额外软件。

##### 测试健康检查

```powershell
# 健康检查
Invoke-RestMethod -Uri "http://localhost:8080/v1/healthz" -Method Get

# 预期输出:
# status
# ------
# ok
```

##### 测试就绪检查

```powershell
# 就绪检查
Invoke-RestMethod -Uri "http://localhost:8080/v1/readyz" -Method Get

# 预期输出:
# ready
# -----
# true
```

##### 测试用户认证

```powershell
# 定义请求体
$body = @{
    username = "john.doe"
    password = "password123"
} | ConvertTo-Json

# 发送认证请求
$response = Invoke-RestMethod `
    -Uri "http://localhost:8080/v1/auth" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

# 显示响应
$response | ConvertTo-Json -Depth 10
```

#### 工具 2: curl

如果安装了 curl（Windows 10 1803+ 内置）：

```powershell
# 健康检查
curl http://localhost:8080/v1/healthz

# 就绪检查
curl http://localhost:8080/v1/readyz

# 用户认证
curl -X POST http://localhost:8080/v1/auth `
  -H "Content-Type: application/json" `
  -d '{"username":"john.doe","password":"password123"}'
```

#### 工具 3: Postman

1. 下载并安装 [Postman](https://www.postman.com/downloads/)
2. 创建新的 Request
3. 设置请求方法和 URL
4. 添加请求体（JSON 格式）
5. 发送请求

### 完整测试脚本

创建 `test-api.ps1` 脚本：

```powershell
# 测试 LDAP 微服务 API

$baseUrl = "http://localhost:8080"

Write-Host "========================================" -ForegroundColor Green
Write-Host "LDAP 微服务 API 测试" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# 测试 1: 健康检查
Write-Host "`n[测试 1] 健康检查 (GET /v1/healthz)" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/v1/healthz" -Method Get
    Write-Host "✓ 成功: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "✗ 失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试 2: 就绪检查
Write-Host "`n[测试 2] 就绪检查 (GET /v1/readyz)" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/v1/readyz" -Method Get
    Write-Host "✓ 成功: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "✗ 失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试 3: 用户认证 (成功)
Write-Host "`n[测试 3] 用户认证 - 有效凭证 (POST /v1/auth)" -ForegroundColor Yellow
$body = @{
    username = "john.doe"
    password = "password123"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod `
        -Uri "$baseUrl/v1/auth" `
        -Method Post `
        -ContentType "application/json" `
        -Body $body
    Write-Host "✓ 成功: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "✗ 失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试 4: 用户认证 (失败)
Write-Host "`n[测试 4] 用户认证 - 无效凭证 (POST /v1/auth)" -ForegroundColor Yellow
$body = @{
    username = "invalid.user"
    password = "wrongpassword"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod `
        -Uri "$baseUrl/v1/auth" `
        -Method Post `
        -ContentType "application/json" `
        -Body $body
    Write-Host "✓ 成功: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "✗ 预期失败: $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "测试完成" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
```

运行测试脚本：

```powershell
.\test-api.ps1
```

## 常见问题

### Q1: 如何验证 .env 文件是否被正确加载？

**A:** 启动服务时查看日志输出。如果 .env 文件被成功加载，日志中会显示配置信息。

```powershell
# 启动服务并查看日志
.\ldap-microservice.exe
```

### Q2: 如何在 PowerShell 中查看已设置的环境变量？

**A:** 使用以下命令：

```powershell
# 查看特定环境变量
$env:LDAP_URL

# 查看所有环境变量
Get-ChildItem env: | Where-Object { $_.Name -like "LDAP*" }

# 或使用 dir 命令
dir env:LDAP*
```

### Q3: 如何清除已设置的环境变量？

**A:** 在 PowerShell 中使用：

```powershell
# 清除单个环境变量
Remove-Item env:LDAP_URL

# 清除所有 LDAP 相关的环境变量
Get-ChildItem env: | Where-Object { $_.Name -like "LDAP*" } | ForEach-Object { Remove-Item "env:$($_.Name)" }
```

### Q4: 服务无法连接到 LDAP 服务器怎么办？

**A:** 检查以下几点：

1. **验证 LDAP_URL 是否正确**
   ```powershell
   $env:LDAP_URL
   ```

2. **检查网络连接**
   ```powershell
   Test-NetConnection -ComputerName ldap.example.com -Port 389
   ```

3. **查看服务日志**
   - 启动服务时查看错误信息
   - 检查 LDAP 服务器是否在线

4. **验证凭证**
   - 确保 LDAP_BIND_DN 和 LDAP_BIND_PASSWORD 正确
   - 尝试使用 LDAP 客户端工具验证凭证

### Q5: 如何在后台运行服务？

**A:** 使用 PowerShell 的 Start-Process 命令：

```powershell
# 在后台启动服务
Start-Process -FilePath ".\ldap-microservice.exe" -WindowStyle Hidden

# 或使用 nohup（如果安装了 Git Bash）
nohup ./ldap-microservice.exe &
```

### Q6: 如何停止运行中的服务？

**A:** 使用 PowerShell 的 Stop-Process 命令：

```powershell
# 查找进程
Get-Process ldap-microservice

# 停止进程
Stop-Process -Name ldap-microservice -Force
```

## 快速开始

### 最快的测试方式

1. **复制 .env 文件**
   ```powershell
   Copy-Item .env.example .env
   ```

2. **编辑 .env 文件**
   - 打开 `.env` 文件
   - 修改 LDAP 连接参数

3. **编译并运行**
   ```powershell
   go build -o ldap-microservice.exe
   .\ldap-microservice.exe
   ```

4. **测试 API**
   ```powershell
   Invoke-RestMethod -Uri "http://localhost:8080/v1/healthz" -Method Get
   ```

## 相关资源

- [Go 官方文档](https://golang.org/doc/)
- [LDAP 协议文档](https://tools.ietf.org/html/rfc4511)
- [PowerShell 文档](https://docs.microsoft.com/en-us/powershell/)
- [curl 文档](https://curl.se/docs/)

