package routeros

import (
	"context"
)

// Log handles RouterOS log operations
type Log struct {
	client *Client
}

// NewLog creates a new Log resource handler
func NewLog(client *Client) *Log {
	return &Log{client: client}
}

// Print returns all log entries
func (l *Log) Print() (*ResourceList, error) {
	return l.PrintContext(context.Background())
}

// PrintContext returns all log entries with context
func (l *Log) PrintContext(ctx context.Context) (*ResourceList, error) {
	responses, err := l.client.SendContext(ctx, "/log/print")
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// GetByTopics returns log entries matching topics
func (l *Log) GetByTopics(topics string) (*ResourceList, error) {
	return l.GetByTopicsContext(context.Background(), topics)
}

// GetByTopicsContext returns log entries matching topics with context
func (l *Log) GetByTopicsContext(ctx context.Context, topics string) (*ResourceList, error) {
	responses, err := l.client.SendContext(ctx, "/log/print", "?topics="+topics)
	if err != nil {
		return nil, err
	}
	return NewResourceList(responses), nil
}

// Warning returns warning level log entries
func (l *Log) Warning() (*ResourceList, error) {
	return l.WarningContext(context.Background())
}

// WarningContext returns warning level log entries with context
func (l *Log) WarningContext(ctx context.Context) (*ResourceList, error) {
	return l.GetByTopicsContext(ctx, "warning")
}

// Error returns error level log entries
func (l *Log) Error() (*ResourceList, error) {
	return l.ErrorContext(context.Background())
}

// ErrorContext returns error level log entries with context
func (l *Log) ErrorContext(ctx context.Context) (*ResourceList, error) {
	return l.GetByTopicsContext(ctx, "error")
}

// Critical returns critical level log entries
func (l *Log) Critical() (*ResourceList, error) {
	return l.CriticalContext(context.Background())
}

// CriticalContext returns critical level log entries with context
func (l *Log) CriticalContext(ctx context.Context) (*ResourceList, error) {
	return l.GetByTopicsContext(ctx, "critical")
}

// Info returns info level log entries
func (l *Log) Info() (*ResourceList, error) {
	return l.InfoContext(context.Background())
}

// InfoContext returns info level log entries with context
func (l *Log) InfoContext(ctx context.Context) (*ResourceList, error) {
	return l.GetByTopicsContext(ctx, "info")
}

// Debug returns debug level log entries
func (l *Log) Debug() (*ResourceList, error) {
	return l.DebugContext(context.Background())
}

// DebugContext returns debug level log entries with context
func (l *Log) DebugContext(ctx context.Context) (*ResourceList, error) {
	return l.GetByTopicsContext(ctx, "debug")
}

// Settings handles log settings operations
type LogSettings struct {
	client *Client
}

// NewLogSettings creates a new LogSettings handler
func NewLogSettings(client *Client) *LogSettings {
	return &LogSettings{client: client}
}

// Get returns log settings
func (s *LogSettings) Get() (map[string]string, error) {
	return s.GetContext(context.Background())
}

// GetContext returns log settings with context
func (s *LogSettings) GetContext(ctx context.Context) (map[string]string, error) {
	responses, err := s.client.SendContext(ctx, "/system/logging/print")
	if err != nil {
		return nil, err
	}
	parser := NewResponseParser(responses)
	return parser.FirstData(), nil
}

// Action handles log action operations
type LogAction struct {
	*BaseResource
}

// NewLogAction creates a new LogAction handler
func NewLogAction(client *Client) *LogAction {
	return &LogAction{
		BaseResource: NewBaseResource(client, "/system/logging/action"),
	}
}

// Rule handles log rule operations
type LogRule struct {
	*BaseResource
}

// NewLogRule creates a new LogRule handler
func NewLogRule(client *Client) *LogRule {
	return &LogRule{
		BaseResource: NewBaseResource(client, "/system/logging"),
	}
}
