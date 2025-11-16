package main

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original env vars
	originalEnv := map[string]string{
		"SERVICE_PORT":              os.Getenv("SERVICE_PORT"),
		"LDAP_URL":                  os.Getenv("LDAP_URL"),
		"LDAP_BIND_DN":              os.Getenv("LDAP_BIND_DN"),
		"LDAP_BIND_PASSWORD":        os.Getenv("LDAP_BIND_PASSWORD"),
		"LDAP_USER_BASE":            os.Getenv("LDAP_USER_BASE"),
		"LDAP_USER_FILTER":          os.Getenv("LDAP_USER_FILTER"),
		"LDAP_USER_DN_ATTR":         os.Getenv("LDAP_USER_DN_ATTR"),
		"LDAP_USE_LDAPS":            os.Getenv("LDAP_USE_LDAPS"),
		"LDAP_USE_STARTTLS":         os.Getenv("LDAP_USE_STARTTLS"),
		"LDAP_INSECURE_SKIP_VERIFY": os.Getenv("LDAP_INSECURE_SKIP_VERIFY"),
	}
	defer func() {
		// Restore original env vars
		for k, v := range originalEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	// Clear all env vars
	for k := range originalEnv {
		os.Unsetenv(k)
	}

	// Test default values
	cfg := LoadConfigFromEnv()
	if cfg.ServicePort != "8080" {
		t.Errorf("expected ServicePort=8080, got %s", cfg.ServicePort)
	}
	if cfg.LDAPURL != "ldap://ldap.example.com:389" {
		t.Errorf("expected LDAPURL=ldap://ldap.example.com:389, got %s", cfg.LDAPURL)
	}
	if cfg.ConnTimeout != 5*time.Second {
		t.Errorf("expected ConnTimeout=5s, got %v", cfg.ConnTimeout)
	}
	if cfg.RequestTimeout != 8*time.Second {
		t.Errorf("expected RequestTimeout=8s, got %v", cfg.RequestTimeout)
	}
	if cfg.UseLDAPS {
		t.Error("expected UseLDAPS=false")
	}
	if cfg.UseStartTLS {
		t.Error("expected UseStartTLS=false")
	}
	if cfg.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify=false")
	}
}

func TestLoadConfigFromEnvWithOverrides(t *testing.T) {
	// Save original env vars
	originalEnv := map[string]string{
		"SERVICE_PORT":              os.Getenv("SERVICE_PORT"),
		"LDAP_URL":                  os.Getenv("LDAP_URL"),
		"LDAP_USE_LDAPS":            os.Getenv("LDAP_USE_LDAPS"),
		"LDAP_USE_STARTTLS":         os.Getenv("LDAP_USE_STARTTLS"),
		"LDAP_INSECURE_SKIP_VERIFY": os.Getenv("LDAP_INSECURE_SKIP_VERIFY"),
	}
	defer func() {
		for k, v := range originalEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	// Set custom env vars
	os.Setenv("SERVICE_PORT", "9090")
	os.Setenv("LDAP_URL", "ldaps://custom.ldap.com:636")
	os.Setenv("LDAP_USE_LDAPS", "1")
	os.Setenv("LDAP_INSECURE_SKIP_VERIFY", "1")

	cfg := LoadConfigFromEnv()
	if cfg.ServicePort != "9090" {
		t.Errorf("expected ServicePort=9090, got %s", cfg.ServicePort)
	}
	if cfg.LDAPURL != "ldaps://custom.ldap.com:636" {
		t.Errorf("expected LDAPURL=ldaps://custom.ldap.com:636, got %s", cfg.LDAPURL)
	}
	if !cfg.UseLDAPS {
		t.Error("expected UseLDAPS=true")
	}
	if !cfg.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify=true")
	}
}

func TestConfigToMap(t *testing.T) {
	cfg := &Config{
		ServicePort:       "8080",
		LDAPURL:           "ldap://localhost:389",
		UserSearchBase:    "dc=example,dc=com",
		UserSearchFilter:  "(uid=%s)",
		UseLDAPS:          false,
		UseStartTLS:       false,
	}

	m := cfg.ToMap()
	if m["ServicePort"] != "8080" {
		t.Errorf("expected ServicePort=8080 in map, got %v", m["ServicePort"])
	}
	if m["LDAPURL"] != "ldap://localhost:389" {
		t.Errorf("expected LDAPURL=ldap://localhost:389 in map, got %v", m["LDAPURL"])
	}
}

func TestGetEnv(t *testing.T) {
	// Save original
	original := os.Getenv("TEST_VAR")
	defer func() {
		if original == "" {
			os.Unsetenv("TEST_VAR")
		} else {
			os.Setenv("TEST_VAR", original)
		}
	}()

	// Test with env var set
	os.Setenv("TEST_VAR", "custom_value")
	result := getEnv("TEST_VAR", "default_value")
	if result != "custom_value" {
		t.Errorf("expected custom_value, got %s", result)
	}

	// Test with env var not set
	os.Unsetenv("TEST_VAR")
	result = getEnv("TEST_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("expected default_value, got %s", result)
	}
}

