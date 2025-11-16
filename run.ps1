# LDAP 微服务启动脚本
# 用于在 Windows 环境下启动 LDAP 微服务

param(
    [switch]$Build = $false,
    [switch]$Test = $false,
    [string]$EnvFile = ".env"
)

# 颜色定义
$colors = @{
    Success = "Green"
    Error   = "Red"
    Warning = "Yellow"
    Info    = "Cyan"
}

# 打印带颜色的消息
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

# 打印分隔线
function Write-Separator {
    Write-Host "========================================" -ForegroundColor Cyan
}

# 检查 .env 文件
function Check-EnvFile {
    if (-not (Test-Path $EnvFile)) {
        Write-ColorOutput "警告: $EnvFile 文件不存在" -Color $colors.Warning
        Write-ColorOutput "正在从 .env.example 创建 $EnvFile..." -Color $colors.Info
        
        if (Test-Path ".env.example") {
            Copy-Item ".env.example" $EnvFile
            Write-ColorOutput "✓ $EnvFile 已创建" -Color $colors.Success
            Write-ColorOutput "请编辑 $EnvFile 文件配置 LDAP 连接参数" -Color $colors.Warning
        } else {
            Write-ColorOutput "✗ 错误: .env.example 文件不存在" -Color $colors.Error
            return $false
        }
    } else {
        Write-ColorOutput "✓ 找到 $EnvFile 文件" -Color $colors.Success
    }
    return $true
}

# 编译项目
function Build-Project {
    Write-ColorOutput "`n正在编译项目..." -Color $colors.Info
    
    try {
        go build -o ldap-microservice.exe
        if ($LASTEXITCODE -eq 0) {
            Write-ColorOutput "✓ 编译成功" -Color $colors.Success
            return $true
        } else {
            Write-ColorOutput "✗ 编译失败" -Color $colors.Error
            return $false
        }
    } catch {
        Write-ColorOutput "✗ 编译错误: $_" -Color $colors.Error
        return $false
    }
}

# 启动服务
function Start-Service {
    Write-ColorOutput "`n正在启动 LDAP 微服务..." -Color $colors.Info
    Write-ColorOutput "服务将监听 http://localhost:8080" -Color $colors.Info
    Write-ColorOutput "按 Ctrl+C 停止服务" -Color $colors.Warning
    Write-Separator
    
    try {
        .\ldap-microservice.exe
    } catch {
        Write-ColorOutput "✗ 启动失败: $_" -Color $colors.Error
        return $false
    }
}

# 运行测试
function Run-Tests {
    Write-ColorOutput "`n正在运行 API 测试..." -Color $colors.Info
    
    if (-not (Test-Path "test-api.ps1")) {
        Write-ColorOutput "✗ 错误: test-api.ps1 文件不存在" -Color $colors.Error
        return $false
    }
    
    try {
        & .\test-api.ps1
        return $LASTEXITCODE -eq 0
    } catch {
        Write-ColorOutput "✗ 测试失败: $_" -Color $colors.Error
        return $false
    }
}

# 显示帮助信息
function Show-Help {
    Write-Separator
    Write-ColorOutput "LDAP 微服务启动脚本" -Color $colors.Info
    Write-Separator
    Write-Host ""
    Write-Host "用法: .\run.ps1 [选项]"
    Write-Host ""
    Write-Host "选项:"
    Write-Host "  -Build    编译项目后启动服务"
    Write-Host "  -Test     启动服务后运行 API 测试"
    Write-Host "  -EnvFile  指定环境变量文件 (默认: .env)"
    Write-Host "  -Help     显示此帮助信息"
    Write-Host ""
    Write-Host "示例:"
    Write-Host "  .\run.ps1                    # 直接启动服务"
    Write-Host "  .\run.ps1 -Build             # 编译后启动"
    Write-Host "  .\run.ps1 -Build -Test       # 编译、启动并测试"
    Write-Host "  .\run.ps1 -EnvFile .env.dev  # 使用自定义环境文件"
    Write-Host ""
    Write-Separator
}

# 主函数
function Main {
    Write-Separator
    Write-ColorOutput "LDAP 微服务启动脚本" -Color $colors.Info
    Write-Separator
    Write-Host ""

    # 检查 .env 文件
    if (-not (Check-EnvFile)) {
        Write-ColorOutput "✗ 无法继续: .env 文件检查失败" -Color $colors.Error
        exit 1
    }

    # 编译项目（如果指定了 -Build）
    if ($Build) {
        if (-not (Build-Project)) {
            Write-ColorOutput "✗ 无法继续: 编译失败" -Color $colors.Error
            exit 1
        }
    }

    # 检查可执行文件是否存在
    if (-not (Test-Path "ldap-microservice.exe")) {
        Write-ColorOutput "✗ 错误: ldap-microservice.exe 不存在" -Color $colors.Error
        Write-ColorOutput "请先编译项目: go build -o ldap-microservice.exe" -Color $colors.Warning
        exit 1
    }

    # 启动服务
    Start-Service

    # 运行测试（如果指定了 -Test）
    if ($Test) {
        Run-Tests
    }
}

# 处理命令行参数
if ($args -contains "-Help" -or $args -contains "-h" -or $args -contains "/?") {
    Show-Help
    exit 0
}

# 运行主函数
Main

