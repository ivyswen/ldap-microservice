# Windows 环境配置支持 - 功能总结

本文档总结了为 LDAP 微服务项目添加的 Windows 环境配置支持功能。

## 新增功能概览

### 1. .env 文件支持 ✅

**文件**: `config.go`

- 添加了 `github.com/joho/godotenv` 依赖
- 修改 `LoadConfigFromEnv()` 函数以自动加载 `.env` 文件
- 如果 `.env` 文件不存在，程序仍然可以从环境变量加载配置

**优势**:
- 无需手动设置环境变量
- 配置集中管理
- 易于在不同环境间切换
- 不会污染系统环境变量

### 2. 配置文件示例

**文件**: `.env.example`

- 包含所有可配置的环境变量
- 提供详细的中英文注释
- 包含 3 个常见配置场景示例：
  - 本地 OpenLDAP 测试环境
  - Active Directory (LDAPS)
  - Active Directory (StartTLS)

### 3. Windows 测试指南

**文件**: `WINDOWS_TESTING_GUIDE.md`

完整的 Windows 环境测试指南，包含：
- 环境变量配置方法（3 种方式）
- .env 文件使用指南
- 服务启动方法
- API 端点测试方法
- 常见问题解答
- 故障排除指南

### 4. 快速开始指南

**文件**: `QUICK_START_WINDOWS.md`

针对 Windows 用户的快速开始指南：
- 5 分钟快速开始步骤
- 常见配置场景
- 完整 API 测试示例
- 环境变量说明表
- 故障排除

### 5. PowerShell 脚本

#### 5.1 启动脚本

**文件**: `run.ps1`

功能：
- 自动检查 .env 文件
- 支持编译项目
- 启动服务
- 运行 API 测试

用法：
```powershell
.\run.ps1                    # 直接启动
.\run.ps1 -Build             # 编译后启动
.\run.ps1 -Build -Test       # 编译、启动并测试
```

#### 5.2 API 测试脚本

**文件**: `test-api.ps1`

功能：
- 测试健康检查端点
- 测试就绪检查端点
- 测试有效认证
- 测试无效认证
- 彩色输出和详细报告

用法：
```powershell
.\test-api.ps1                                    # 使用默认参数
.\test-api.ps1 -BaseUrl "http://localhost:9090"  # 自定义 URL
.\test-api.ps1 -Username "user" -Password "pass" # 自定义凭证
```

#### 5.3 环境设置脚本

**文件**: `setup-env.ps1`

功能：
- 快速设置环境变量
- 支持临时和永久设置
- 自动创建 .env 文件
- 支持多种 LDAP 配置场景

用法：
```powershell
.\setup-env.ps1                                    # 使用默认值
.\setup-env.ps1 -Permanent                        # 永久设置
.\setup-env.ps1 -LdapUrl "ldaps://ad.com:636" \
  -UseLDAPS -BindDN "CN=Admin,CN=Users,DC=com"   # 自定义配置
```

### 6. Makefile

**文件**: `Makefile`

提供常用命令的快捷方式：
- `make build` - 编译项目
- `make run` - 编译并运行
- `make test` - 运行单元测试
- `make test-api` - 运行 API 测试
- `make clean` - 清理编译产物
- `make install-deps` - 安装依赖
- `make fmt` - 格式化代码
- `make lint` - 代码检查

### 7. 文档更新

#### 7.1 README.md

更新内容：
- 添加 .env 文件使用指南
- 添加 Windows 测试部分
- 添加 API 测试示例
- 更新开发指南
- 添加 Windows 开发说明

#### 7.2 .gitignore

添加：
- `.env` - 本地环境配置
- `.env.local` - 本地覆盖配置
- `.env.*.local` - 环境特定配置

## 使用流程

### 最快的开始方式

```powershell
# 1. 复制配置文件
Copy-Item .env.example .env

# 2. 编辑 .env 文件（使用任何文本编辑器）
# 修改 LDAP 连接参数

# 3. 编译项目
go build -o ldap-microservice.exe

# 4. 启动服务
.\ldap-microservice.exe

# 5. 测试 API（在新窗口）
.\test-api.ps1
```

### 使用脚本的方式

```powershell
# 1. 设置环境
.\setup-env.ps1 -LdapUrl "ldap://ldap.example.com:389" -Permanent

# 2. 编译并启动
.\run.ps1 -Build

# 3. 测试 API
.\test-api.ps1
```

## 文件清单

### 新增文件

| 文件 | 类型 | 说明 |
|------|------|------|
| `.env.example` | 配置 | 环境变量示例文件 |
| `WINDOWS_TESTING_GUIDE.md` | 文档 | Windows 测试完整指南 |
| `QUICK_START_WINDOWS.md` | 文档 | Windows 快速开始指南 |
| `WINDOWS_ENV_SETUP_SUMMARY.md` | 文档 | 本文件 |
| `run.ps1` | 脚本 | PowerShell 启动脚本 |
| `test-api.ps1` | 脚本 | PowerShell API 测试脚本 |
| `setup-env.ps1` | 脚本 | PowerShell 环境设置脚本 |
| `Makefile` | 构建 | Make 构建文件 |

### 修改文件

| 文件 | 修改内容 |
|------|---------|
| `config.go` | 添加 godotenv 支持 |
| `go.mod` | 添加 godotenv 依赖 |
| `go.sum` | 更新依赖哈希 |
| `README.md` | 添加 .env 和 Windows 测试说明 |
| `.gitignore` | 添加 .env 文件忽略规则 |

## 依赖管理

### 新增依赖

- `github.com/joho/godotenv` - .env 文件加载库

### 依赖安装

```bash
go get github.com/joho/godotenv
go mod tidy
```

## 配置优先级

程序加载配置的优先级（从高到低）：

1. 环境变量（系统或当前会话设置）
2. .env 文件中的变量
3. 代码中的默认值

这意味着：
- 系统环境变量会覆盖 .env 文件
- .env 文件会覆盖代码默认值
- 如果都没有设置，使用代码默认值

## 常见使用场景

### 场景 1: 开发环境

```powershell
# 使用 .env 文件管理配置
Copy-Item .env.example .env
# 编辑 .env 文件
.\run.ps1 -Build -Test
```

### 场景 2: 测试环境

```powershell
# 使用脚本快速设置
.\setup-env.ps1 -LdapUrl "ldap://test-ldap:389"
go build
.\test-api.ps1
```

### 场景 3: 生产环境

```powershell
# 使用系统环境变量
[Environment]::SetEnvironmentVariable("LDAP_URL", "ldaps://prod-ldap:636", "Machine")
# 重启应用
.\ldap-microservice.exe
```

## 故障排除

### 问题: PowerShell 脚本无法执行

**解决方案**:
```powershell
# 允许执行本地脚本
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 问题: .env 文件未被加载

**检查**:
1. 确保 .env 文件在项目根目录
2. 查看服务启动日志
3. 验证 .env 文件格式正确

### 问题: 环境变量优先级问题

**解决方案**:
- 清除系统环境变量：`Remove-Item env:LDAP_URL`
- 或在 .env 文件中设置

## 下一步

1. 查看 [QUICK_START_WINDOWS.md](QUICK_START_WINDOWS.md) 快速开始
2. 查看 [WINDOWS_TESTING_GUIDE.md](WINDOWS_TESTING_GUIDE.md) 完整指南
3. 查看 [README.md](README.md) 了解更多功能
4. 查看 [.env.example](.env.example) 了解所有配置选项

## 相关资源

- [godotenv 库文档](https://github.com/joho/godotenv)
- [PowerShell 文档](https://docs.microsoft.com/en-us/powershell/)
- [Go 官方文档](https://golang.org/doc/)
- [LDAP 协议](https://tools.ietf.org/html/rfc4511)

