package main

import (
	"testing"
	"time"
)

func TestLDAPClientClose(t *testing.T) {
	// Test Close with nil connection
	client := &LDAPClient{
		cfg:  &Config{},
		conn: nil,
	}
	// Should not panic
	client.Close()
}

func TestNewLDAPClientInvalidURL(t *testing.T) {
	cfg := &Config{
		LDAPURL:     "invalid://url",
		ConnTimeout: 2 * time.Second,
	}

	client, err := NewLDAPClient(cfg)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
	if client != nil {
		t.Error("expected nil client for invalid URL")
	}

	// Check if it's an LDAPError
	if !IsLDAPError(err) {
		t.Errorf("expected LDAPError, got %T", err)
	}

	// Check error code
	code := GetLDAPErrorCode(err)
	if code != ErrConnectionFailed {
		t.Errorf("expected ErrConnectionFailed, got %v", code)
	}
}

func TestNewLDAPClientConnectionTimeout(t *testing.T) {
	// Use an unreachable address to trigger timeout
	cfg := &Config{
		LDAPURL:     "ldap://192.0.2.1:389", // TEST-NET-1 (unreachable)
		ConnTimeout: 100 * time.Millisecond,  // Very short timeout
	}

	client, err := NewLDAPClient(cfg)
	if err == nil {
		t.Error("expected error for connection timeout")
		if client != nil {
			client.Close()
		}
	}
}

func TestIsLDAPError(t *testing.T) {
	ldapErr := NewLDAPError(ErrConnectionFailed, "test error")
	if !IsLDAPError(ldapErr) {
		t.Error("expected IsLDAPError to return true for LDAPError")
	}

	regularErr := NewLDAPError(ErrConnectionFailed, "test")
	if !IsLDAPError(regularErr) {
		t.Error("expected IsLDAPError to return true")
	}
}

func TestGetLDAPErrorCode(t *testing.T) {
	ldapErr := NewLDAPError(ErrUserNotFound, "user not found")
	code := GetLDAPErrorCode(ldapErr)
	if code != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", code)
	}

	// Test with non-LDAPError
	code = GetLDAPErrorCode(nil)
	if code != "" {
		t.Errorf("expected empty code for nil error, got %v", code)
	}
}

func TestLDAPErrorString(t *testing.T) {
	// Test without cause
	err := NewLDAPError(ErrConnectionFailed, "connection failed")
	errStr := err.Error()
	if errStr != "[connection_failed] connection failed" {
		t.Errorf("unexpected error string: %s", errStr)
	}

	// Test with cause
	cause := NewLDAPError(ErrBindFailed, "bind failed")
	err = NewLDAPErrorWithCause(ErrConnectionFailed, "connection failed", cause)
	errStr = err.Error()
	if errStr != "[connection_failed] connection failed: [bind_failed] bind failed" {
		t.Errorf("unexpected error string with cause: %s", errStr)
	}
}

func TestLDAPErrorUnwrap(t *testing.T) {
	cause := NewLDAPError(ErrBindFailed, "bind failed")
	err := NewLDAPErrorWithCause(ErrConnectionFailed, "connection failed", cause)

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Error("expected Unwrap to return the cause")
	}
}

