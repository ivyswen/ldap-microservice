# Windows 快速开始指南

本指南帮助你在 Windows 环境下快速启动和测试 LDAP 微服务。

## 前置条件

- Windows 10 或更高版本
- Go 1.25 或更高版本（[下载](https://golang.org/dl/)）
- PowerShell 5.0 或更高版本（Windows 10 内置）
- 可选：curl（Windows 10 1803+ 内置）或 Postman

## 5 分钟快速开始

### 步骤 1: 准备配置文件

```powershell
# 进入项目目录
cd d:\GolandProjects\ldap-microservice

# 复制示例配置文件
Copy-Item .env.example .env
```

### 步骤 2: 编辑配置文件

使用任何文本编辑器打开 `.env` 文件，修改 LDAP 连接参数：

```env
# 最小配置示例
LDAP_URL=ldap://ldap.example.com:389
LDAP_BIND_DN=cn=admin,dc=example,dc=com
LDAP_BIND_PASSWORD=password123
LDAP_USER_BASE=dc=example,dc=com
LDAP_USER_FILTER=(uid=%s)
SERVICE_PORT=8080
```

### 步骤 3: 编译项目

```powershell
go build -o ldap-microservice.exe
```

### 步骤 4: 启动服务

```powershell
.\ldap-microservice.exe
```

你应该看到类似的输出：
```
{"level":"info","time":"2024-01-15T10:30:00Z","message":"Starting LDAP microservice on :8080"}
```

### 步骤 5: 测试 API

在新的 PowerShell 窗口中运行：

```powershell
# 健康检查
Invoke-RestMethod -Uri "http://localhost:8080/v1/healthz" -Method Get

# 预期输出: status = ok
```

## 使用 PowerShell 脚本

### 使用启动脚本

```powershell
# 直接启动（需要已编译）
.\run.ps1

# 编译后启动
.\run.ps1 -Build

# 编译、启动并测试
.\run.ps1 -Build -Test
```

### 使用测试脚本

```powershell
# 运行所有 API 测试
.\test-api.ps1

# 指定自定义 URL
.\test-api.ps1 -BaseUrl "http://localhost:9090"

# 指定自定义用户名和密码
.\test-api.ps1 -Username "testuser" -Password "testpass"
```

## 常见配置场景

### 场景 1: 本地 OpenLDAP 测试

编辑 `.env` 文件：

```env
LDAP_URL=ldap://localhost:389
LDAP_BIND_DN=cn=admin,dc=example,dc=com
LDAP_BIND_PASSWORD=admin
LDAP_USER_BASE=ou=users,dc=example,dc=com
LDAP_USER_FILTER=(uid=%s)
SERVICE_PORT=8080
```

### 场景 2: Active Directory (LDAPS)

编辑 `.env` 文件：

```env
LDAP_URL=ldaps://ad.example.com:636
LDAP_USE_LDAPS=1
LDAP_INSECURE_SKIP_VERIFY=0
LDAP_BIND_DN=CN=ServiceAccount,CN=Users,DC=example,DC=com
LDAP_BIND_PASSWORD=ServicePassword123
LDAP_USER_BASE=CN=Users,DC=example,DC=com
LDAP_USER_FILTER=(sAMAccountName=%s)
SERVICE_PORT=8080
```

### 场景 3: Active Directory (StartTLS)

编辑 `.env` 文件：

```env
LDAP_URL=ldap://ad.example.com:389
LDAP_USE_STARTTLS=1
LDAP_INSECURE_SKIP_VERIFY=0
LDAP_BIND_DN=CN=ServiceAccount,CN=Users,DC=example,DC=com
LDAP_BIND_PASSWORD=ServicePassword123
LDAP_USER_BASE=CN=Users,DC=example,DC=com
LDAP_USER_FILTER=(sAMAccountName=%s)
SERVICE_PORT=8080
```

## 完整 API 测试示例

### 使用 PowerShell

```powershell
# 1. 健康检查
Write-Host "=== 健康检查 ===" -ForegroundColor Green
Invoke-RestMethod -Uri "http://localhost:8080/v1/healthz" -Method Get

# 2. 就绪检查
Write-Host "`n=== 就绪检查 ===" -ForegroundColor Green
Invoke-RestMethod -Uri "http://localhost:8080/v1/readyz" -Method Get

# 3. 用户认证
Write-Host "`n=== 用户认证 ===" -ForegroundColor Green
$body = @{
    username = "john.doe"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod `
    -Uri "http://localhost:8080/v1/auth" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

$response | ConvertTo-Json -Depth 10
```

### 使用 curl

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

## 环境变量说明

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SERVICE_PORT` | 8080 | HTTP 服务端口 |
| `LDAP_URL` | ldap://ldap.example.com:389 | LDAP 服务器 URL |
| `LDAP_BIND_DN` | - | 服务账户 DN |
| `LDAP_BIND_PASSWORD` | - | 服务账户密码 |
| `LDAP_USER_BASE` | dc=example,dc=com | 用户搜索基础 DN |
| `LDAP_USER_FILTER` | (uid=%s) | 用户搜索过滤器 |
| `LDAP_USE_LDAPS` | 0 | 是否使用 LDAPS |
| `LDAP_USE_STARTTLS` | 0 | 是否使用 StartTLS |
| `LDAP_INSECURE_SKIP_VERIFY` | 0 | 是否跳过 TLS 验证 |

## 故障排除

### 问题 1: "无法连接到 LDAP 服务器"

**解决方案：**

1. 验证 LDAP_URL 是否正确
2. 检查网络连接
   ```powershell
   Test-NetConnection -ComputerName ldap.example.com -Port 389
   ```
3. 验证防火墙设置

### 问题 2: "认证失败"

**解决方案：**

1. 验证 LDAP_BIND_DN 和 LDAP_BIND_PASSWORD
2. 检查用户是否存在于 LDAP_USER_BASE
3. 验证 LDAP_USER_FILTER 是否正确

### 问题 3: "TLS 证书错误"

**解决方案：**

1. 对于自签名证书，设置 `LDAP_INSECURE_SKIP_VERIFY=1`（仅用于测试）
2. 或者将证书添加到系统信任存储

### 问题 4: "端口已被占用"

**解决方案：**

```powershell
# 查找占用端口的进程
Get-NetTCPConnection -LocalPort 8080

# 修改 SERVICE_PORT 为其他端口
# 编辑 .env 文件，改为 SERVICE_PORT=9090
```

## 下一步

- 查看完整的 [WINDOWS_TESTING_GUIDE.md](WINDOWS_TESTING_GUIDE.md)
- 查看 [README.md](README.md) 了解更多功能
- 查看 [.env.example](.env.example) 了解所有配置选项

## 获取帮助

- 查看项目文档
- 检查服务日志输出
- 运行 `.\test-api.ps1` 进行诊断

## 相关资源

- [Go 官方文档](https://golang.org/doc/)
- [LDAP 协议](https://tools.ietf.org/html/rfc4511)
- [PowerShell 文档](https://docs.microsoft.com/en-us/powershell/)

