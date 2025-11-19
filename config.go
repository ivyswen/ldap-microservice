package main

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServicePort        string
	LDAPURL            string
	BindDN             string
	BindPassword       string
	UserSearchBase     string
	UserSearchFilter   string // e.g. "(uid=%s)" or "(sAMAccountName=%s)"
	UserDNAttr         string // optional
	ReturnAttributes   []string
	ConnTimeout        time.Duration
	RequestTimeout     time.Duration
	UseLDAPS           bool
	UseStartTLS        bool
	InsecureSkipVerify bool
	BasePath           string // URL 路径前缀，例如 "/api" 或 "/ldap"
	LogLevel           string // 日志级别: debug, info, warn, error
	LogFile            string // 日志文件路径，为空则只输出到控制台
}

func LoadConfigFromEnv() *Config {
	// 尝试从 .env 文件加载环境变量（如果存在）
	// 如果 .env 文件不存在，godotenv.Load() 会返回错误但不会中断程序
	_ = godotenv.Load()

	c := &Config{
		ServicePort:        getEnv("SERVICE_PORT", "8080"),
		LDAPURL:            getEnv("LDAP_URL", "ldap://ldap.example.com:389"),
		BindDN:             os.Getenv("LDAP_BIND_DN"),
		BindPassword:       os.Getenv("LDAP_BIND_PASSWORD"),
		UserSearchBase:     getEnv("LDAP_USER_BASE", "dc=example,dc=com"),
		UserSearchFilter:   getEnv("LDAP_USER_FILTER", "(uid=%s)"),
		UserDNAttr:         os.Getenv("LDAP_USER_DN_ATTR"),
		ReturnAttributes:   []string{"cn", "mail", "uid"},
		ConnTimeout:        5 * time.Second,
		RequestTimeout:     8 * time.Second,
		UseLDAPS:           getEnv("LDAP_USE_LDAPS", "") == "1",
		UseStartTLS:        getEnv("LDAP_USE_STARTTLS", "") == "1",
		InsecureSkipVerify: getEnv("LDAP_INSECURE_SKIP_VERIFY", "") == "1",
		BasePath:           normalizePath(getEnv("BASE_PATH", "")),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFile:            getEnv("LOG_FILE", "app.log"),
	}
	return c
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// normalizePath 规范化 URL 路径前缀
// - 移除末尾的 /
// - 确保开头有 /（如果路径非空）
// - 空字符串保持为空
func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	// 移除末尾的 /
	path = strings.TrimSuffix(path, "/")
	// 确保开头有 /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

func (c *Config) ToMap() map[string]any {
	return map[string]any{
		"ServicePort":      c.ServicePort,
		"LDAPURL":          c.LDAPURL,
		"UserSearchBase":   c.UserSearchBase,
		"UserSearchFilter": c.UserSearchFilter,
		"UseLDAPS":         c.UseLDAPS,
		"UseStartTLS":      c.UseStartTLS,
		"BasePath":         c.BasePath,
		"LogLevel":         c.LogLevel,
		"LogFile":          c.LogFile,
	}
}
