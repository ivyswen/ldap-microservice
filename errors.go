package main

import "fmt"

// LDAPErrorCode represents different types of LDAP errors
type LDAPErrorCode string

const (
	// Connection errors
	ErrConnectionFailed LDAPErrorCode = "connection_failed"
	ErrConnectionTimeout LDAPErrorCode = "connection_timeout"
	ErrTLSFailed LDAPErrorCode = "tls_failed"

	// Authentication errors
	ErrBindFailed LDAPErrorCode = "bind_failed"
	ErrInvalidCredentials LDAPErrorCode = "invalid_credentials"

	// Search errors
	ErrSearchFailed LDAPErrorCode = "search_failed"
	ErrUserNotFound LDAPErrorCode = "user_not_found"
	ErrSearchTimeout LDAPErrorCode = "search_timeout"

	// Configuration errors
	ErrInvalidConfig LDAPErrorCode = "invalid_config"
)

// LDAPError represents a detailed LDAP error with code, message, and underlying error
type LDAPError struct {
	Code  LDAPErrorCode
	Msg   string
	Cause error
}

// Error implements the error interface
func (e *LDAPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Msg, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Msg)
}

// Unwrap returns the underlying error for error chain inspection
func (e *LDAPError) Unwrap() error {
	return e.Cause
}

// NewLDAPError creates a new LDAP error with the given code and message
func NewLDAPError(code LDAPErrorCode, msg string) *LDAPError {
	return &LDAPError{
		Code: code,
		Msg:  msg,
	}
}

// NewLDAPErrorWithCause creates a new LDAP error with an underlying cause
func NewLDAPErrorWithCause(code LDAPErrorCode, msg string, cause error) *LDAPError {
	return &LDAPError{
		Code:  code,
		Msg:   msg,
		Cause: cause,
	}
}

// IsLDAPError checks if an error is an LDAPError
func IsLDAPError(err error) bool {
	_, ok := err.(*LDAPError)
	return ok
}

// GetLDAPErrorCode extracts the error code from an error if it's an LDAPError
func GetLDAPErrorCode(err error) LDAPErrorCode {
	if ldapErr, ok := err.(*LDAPError); ok {
		return ldapErr.Code
	}
	return ""
}

