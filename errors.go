// Package routeros provides a comprehensive Go client for MikroTik RouterOS API.
// It supports all RouterOS versions and provides a fluent interface for building queries.
package routeros

import (
	"errors"
	"fmt"
)

// Error definitions for the RouterOS API client
var (
	// ErrNotConnected indicates that the client is not connected to the router
	ErrNotConnected = errors.New("not connected to router")

	// ErrAlreadyConnected indicates that the client is already connected
	ErrAlreadyConnected = errors.New("already connected to router")

	// ErrAuthenticationFailed indicates authentication failure
	ErrAuthenticationFailed = errors.New("authentication failed")

	// ErrInvalidResponse indicates an invalid response from the router
	ErrInvalidResponse = errors.New("invalid response from router")

	// ErrTimeout indicates a timeout occurred
	ErrTimeout = errors.New("operation timed out")

	// ErrConnectionClosed indicates the connection was closed
	ErrConnectionClosed = errors.New("connection closed")

	// ErrInvalidAddress indicates an invalid address format
	ErrInvalidAddress = errors.New("invalid address format")

	// ErrInvalidPort indicates an invalid port number
	ErrInvalidPort = errors.New("invalid port number")

	// ErrInvalidCredentials indicates invalid username or password
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrCommandFailed indicates a command execution failure
	ErrCommandFailed = errors.New("command execution failed")

	// ErrInvalidQuery indicates an invalid query syntax
	ErrInvalidQuery = errors.New("invalid query syntax")

	// ErrResourceNotFound indicates the requested resource was not found
	ErrResourceNotFound = errors.New("resource not found")

	// ErrPermissionDenied indicates insufficient permissions
	ErrPermissionDenied = errors.New("permission denied")

	// ErrTLSError indicates a TLS/SSL related error
	ErrTLSError = errors.New("TLS error")

	// ErrUnsupportedVersion indicates an unsupported RouterOS version
	ErrUnsupportedVersion = errors.New("unsupported RouterOS version")

	// ErrInvalidParameter indicates an invalid parameter was provided
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrEmptyResponse indicates an empty response from the router
	ErrEmptyResponse = errors.New("empty response from router")

	// ErrTrap indicates a trap response from the router (error condition)
	ErrTrap = errors.New("trap received from router")

	// ErrFatal indicates a fatal error from the router
	ErrFatal = errors.New("fatal error from router")

	// ErrInterrupted indicates the operation was interrupted
	ErrInterrupted = errors.New("operation interrupted")
)

// RouterError represents an error response from the RouterOS device
type RouterError struct {
	Category int               // Error category (optional)
	Message  string            // Error message from the router
	Tag      string            // Command tag (if any)
	Fields   map[string]string // Additional error fields
}

// Error implements the error interface
func (e *RouterError) Error() string {
	if e.Category > 0 {
		return fmt.Sprintf("router error (category %d): %s", e.Category, e.Message)
	}
	return fmt.Sprintf("router error: %s", e.Message)
}

// Is implements errors.Is interface for comparison
func (e *RouterError) Is(target error) bool {
	switch e.Message {
	case "invalid user name or password":
		return target == ErrAuthenticationFailed
	case "no such item":
		return target == ErrResourceNotFound
	case "no permissions":
		return target == ErrPermissionDenied
	}
	return false
}

// NewRouterError creates a new RouterError from response fields
func NewRouterError(fields map[string]string) *RouterError {
	err := &RouterError{
		Fields: fields,
	}
	if msg, ok := fields["message"]; ok {
		err.Message = msg
	}
	if tag, ok := fields["tag"]; ok {
		err.Tag = tag
	}
	if cat, ok := fields["category"]; ok {
		fmt.Sscanf(cat, "%d", &err.Category)
	}
	return err
}

// ConnectionError represents a connection-related error
type ConnectionError struct {
	Addr     string // Connection address
	Port     int    // Connection port
	Internal error  // Underlying error
}

// Error implements the error interface
func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error to %s:%d: %v", e.Addr, e.Port, e.Internal)
}

// Unwrap returns the underlying error
func (e *ConnectionError) Unwrap() error {
	return e.Internal
}

// NewConnectionError creates a new ConnectionError
func NewConnectionError(addr string, port int, err error) *ConnectionError {
	return &ConnectionError{
		Addr:     addr,
		Port:     port,
		Internal: err,
	}
}

// CommandError represents a command execution error
type CommandError struct {
	Command string            // The command that failed
	Args    map[string]string // Command arguments
	Err     error             // Underlying error
}

// Error implements the error interface
func (e *CommandError) Error() string {
	return fmt.Sprintf("command '%s' failed: %v", e.Command, e.Err)
}

// Unwrap returns the underlying error
func (e *CommandError) Unwrap() error {
	return e.Err
}

// NewCommandError creates a new CommandError
func NewCommandError(command string, args map[string]string, err error) *CommandError {
	return &CommandError{
		Command: command,
		Args:    args,
		Err:     err,
	}
}

// QueryError represents a query building error
type QueryError struct {
	Query string // The problematic query part
	Err   error  // Underlying error
}

// Error implements the error interface
func (e *QueryError) Error() string {
	return fmt.Sprintf("query error at '%s': %v", e.Query, e.Err)
}

// Unwrap returns the underlying error
func (e *QueryError) Unwrap() error {
	return e.Err
}

// NewQueryError creates a new QueryError
func NewQueryError(query string, err error) *QueryError {
	return &QueryError{
		Query: query,
		Err:   err,
	}
}

// TimeoutError represents a timeout error with duration information
type TimeoutError struct {
	Operation string        // The operation that timed out
	Duration  string        // The timeout duration
	Err       error         // Underlying error (if any)
}

// Error implements the error interface
func (e *TimeoutError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("timeout during %s after %s: %v", e.Operation, e.Duration, e.Err)
	}
	return fmt.Sprintf("timeout during %s after %s", e.Operation, e.Duration)
}

// Unwrap returns the underlying error
func (e *TimeoutError) Unwrap() error {
	return e.Err
}

// NewTimeoutError creates a new TimeoutError
func NewTimeoutError(operation, duration string, err error) *TimeoutError {
	return &TimeoutError{
		Operation: operation,
		Duration:  duration,
		Err:       err,
	}
}
