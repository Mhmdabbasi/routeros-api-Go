package routeros

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ResponseParser helps parse RouterOS API responses
type ResponseParser struct {
	responses []*Response
}

// NewResponseParser creates a new ResponseParser
func NewResponseParser(responses []*Response) *ResponseParser {
	return &ResponseParser{responses: responses}
}

// First returns the first response or nil
func (p *ResponseParser) First() *Response {
	if len(p.responses) == 0 {
		return nil
	}
	return p.responses[0]
}

// Last returns the last response or nil
func (p *ResponseParser) Last() *Response {
	if len(p.responses) == 0 {
		return nil
	}
	return p.responses[len(p.responses)-1]
}

// All returns all responses
func (p *ResponseParser) All() []*Response {
	return p.responses
}

// Data returns all data from all responses
func (p *ResponseParser) Data() []map[string]string {
	data := make([]map[string]string, 0, len(p.responses))
	for _, resp := range p.responses {
		if resp.Type == ResponseReply {
			data = append(data, resp.Data)
		}
	}
	return data
}

// FirstData returns data from the first reply response
func (p *ResponseParser) FirstData() map[string]string {
	for _, resp := range p.responses {
		if resp.Type == ResponseReply {
			return resp.Data
		}
	}
	return nil
}

// Done checks if command completed successfully
func (p *ResponseParser) Done() bool {
	for _, resp := range p.responses {
		if resp.Type == ResponseDone {
			return true
		}
	}
	return false
}

// Trap checks if a trap occurred
func (p *ResponseParser) Trap() bool {
	for _, resp := range p.responses {
		if resp.Type == ResponseTrap {
			return true
		}
	}
	return false
}

// Fatal checks if a fatal error occurred
func (p *ResponseParser) Fatal() bool {
	for _, resp := range p.responses {
		if resp.Type == ResponseFatal {
			return true
		}
	}
	return false
}

// Error returns the error message if any
func (p *ResponseParser) Error() string {
	for _, resp := range p.responses {
		if resp.Type == ResponseTrap || resp.Type == ResponseFatal {
			return resp.Message
		}
	}
	return ""
}

// Count returns the number of reply responses
func (p *ResponseParser) Count() int {
	count := 0
	for _, resp := range p.responses {
		if resp.Type == ResponseReply {
			count++
		}
	}
	return count
}

// Value returns a specific value from the first data response
func (p *ResponseParser) Value(key string) string {
	data := p.FirstData()
	if data != nil {
		return data[key]
	}
	return ""
}

// Values returns all values for a key across all responses
func (p *ResponseParser) Values(key string) []string {
	values := make([]string, 0)
	for _, resp := range p.responses {
		if resp.Type == ResponseReply {
			if val, ok := resp.Data[key]; ok {
				values = append(values, val)
			}
		}
	}
	return values
}

// ValueInt returns a value as integer
func (p *ResponseParser) ValueInt(key string) (int, error) {
	val := p.Value(key)
	return strconv.Atoi(val)
}

// ValueInt64 returns a value as int64
func (p *ResponseParser) ValueInt64(key string) (int64, error) {
	val := p.Value(key)
	return strconv.ParseInt(val, 10, 64)
}

// ValueBool returns a value as boolean
func (p *ResponseParser) ValueBool(key string) (bool, error) {
	val := p.Value(key)
	val = strings.ToLower(val)
	return val == "true" || val == "yes" || val == "1", nil
}

// ValueFloat returns a value as float64
func (p *ResponseParser) ValueFloat(key string) (float64, error) {
	val := p.Value(key)
	return strconv.ParseFloat(val, 64)
}

// ValueDuration returns a value as duration (RouterOS time format)
func (p *ResponseParser) ValueDuration(key string) (time.Duration, error) {
	val := p.Value(key)
	return ParseDuration(val)
}

// ValueIP returns a value as IP address
func (p *ResponseParser) ValueIP(key string) (string, error) {
	val := p.Value(key)
	// Basic IP validation could be added here
	return val, nil
}

// ToJSON converts responses to JSON
func (p *ResponseParser) ToJSON() ([]byte, error) {
	return json.Marshal(p.Data())
}

// ToJSONIndent converts responses to indented JSON
func (p *ResponseParser) ToJSONIndent() ([]byte, error) {
	return json.MarshalIndent(p.Data(), "", "  ")
}

// ParseDuration parses RouterOS time format to time.Duration
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}

	// Handle special values
	switch s {
	case "none":
		return 0, nil
	case "infinity":
		return time.Duration(1<<63 - 1), nil
	}

	var total time.Duration
	var num int64
	var unit string

	// Parse components like "1d2h3m4s"
	for len(s) > 0 {
		// Find number
		i := 0
		for i < len(s) && (s[i] >= '0' && s[i] <= '9') {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("invalid duration: %s", s)
		}

		_, err := fmt.Sscanf(s[:i], "%d", &num)
		if err != nil {
			return 0, err
		}

		// Find unit
		start := i
		for i < len(s) && (s[i] < '0' || s[i] > '9') {
			i++
		}
		unit = s[start:i]

		// Convert to duration
		switch unit {
		case "d":
			total += time.Duration(num) * 24 * time.Hour
		case "h":
			total += time.Duration(num) * time.Hour
		case "m":
			total += time.Duration(num) * time.Minute
		case "s":
			total += time.Duration(num) * time.Second
		case "ms":
			total += time.Duration(num) * time.Millisecond
		case "us":
			total += time.Duration(num) * time.Microsecond
		default:
			return 0, fmt.Errorf("unknown duration unit: %s", unit)
		}

		s = s[i:]
	}

	return total, nil
}

// ParseBool parses RouterOS boolean values
func ParseBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "yes" || s == "1"
}

// ParseBytes parses RouterOS byte values (e.g., "1KiB", "1MiB")
func ParseBytes(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	var num float64
	var unit string

	_, err := fmt.Sscanf(s, "%f%s", &num, &unit)
	if err != nil {
		return 0, err
	}

	multipliers := map[string]int64{
		"B":   1,
		"KiB": 1024,
		"MiB": 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"TiB": 1024 * 1024 * 1024 * 1024,
		"KB":  1000,
		"MB":  1000 * 1000,
		"GB":  1000 * 1000 * 1000,
		"TB":  1000 * 1000 * 1000 * 1000,
	}

	if mult, ok := multipliers[unit]; ok {
		return int64(num * float64(mult)), nil
	}

	return int64(num), nil
}

// Resource represents a RouterOS resource with an ID
type Resource struct {
	ID     string            // Resource ID (.id)
	Name   string            // Resource name
	Fields map[string]string // All fields
}

// NewResource creates a new Resource from response data
func NewResource(data map[string]string) *Resource {
	r := &Resource{
		ID:     data[".id"],
		Name:   data["name"],
		Fields: data,
	}
	return r
}

// Get returns a field value
func (r *Resource) Get(key string) string {
	return r.Fields[key]
}

// GetInt returns a field value as integer
func (r *Resource) GetInt(key string) (int, error) {
	return strconv.Atoi(r.Fields[key])
}

// GetBool returns a field value as boolean
func (r *Resource) GetBool(key string) bool {
	return ParseBool(r.Fields[key])
}

// ToMap returns all fields as a map
func (r *Resource) ToMap() map[string]string {
	return r.Fields
}

// ToJSON converts the resource to JSON
func (r *Resource) ToJSON() ([]byte, error) {
	return json.Marshal(r.Fields)
}

// ResourceList represents a list of resources
type ResourceList struct {
	items []*Resource
}

// NewResourceList creates a ResourceList from responses
func NewResourceList(responses []*Response) *ResourceList {
	items := make([]*Resource, 0)
	for _, resp := range responses {
		if resp.Type == ResponseReply {
			items = append(items, NewResource(resp.Data))
		}
	}
	return &ResourceList{items: items}
}

// All returns all resources
func (l *ResourceList) All() []*Resource {
	return l.items
}

// First returns the first resource or nil
func (l *ResourceList) First() *Resource {
	if len(l.items) == 0 {
		return nil
	}
	return l.items[0]
}

// Get returns a resource by ID
func (l *ResourceList) Get(id string) *Resource {
	for _, item := range l.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

// FindByName returns a resource by name
func (l *ResourceList) FindByName(name string) *Resource {
	for _, item := range l.items {
		if item.Name == name {
			return item
		}
	}
	return nil
}

// Len returns the number of resources
func (l *ResourceList) Len() int {
	return len(l.items)
}

// Filter filters resources by a predicate function
func (l *ResourceList) Filter(predicate func(*Resource) bool) *ResourceList {
	items := make([]*Resource, 0)
	for _, item := range l.items {
		if predicate(item) {
			items = append(items, item)
		}
	}
	return &ResourceList{items: items}
}

// Map transforms each resource using a function
func (l *ResourceList) Map(transform func(*Resource) interface{}) []interface{} {
	result := make([]interface{}, 0, len(l.items))
	for _, item := range l.items {
		result = append(result, transform(item))
	}
	return result
}

// IDs returns all resource IDs
func (l *ResourceList) IDs() []string {
	ids := make([]string, 0, len(l.items))
	for _, item := range l.items {
		ids = append(ids, item.ID)
	}
	return ids
}

// Names returns all resource names
func (l *ResourceList) Names() []string {
	names := make([]string, 0, len(l.items))
	for _, item := range l.items {
		if item.Name != "" {
			names = append(names, item.Name)
		}
	}
	return names
}

// ToJSON converts the resource list to JSON
func (l *ResourceList) ToJSON() ([]byte, error) {
	data := make([]map[string]string, 0, len(l.items))
	for _, item := range l.items {
		data = append(data, item.Fields)
	}
	return json.Marshal(data)
}

// ToJSONIndent converts the resource list to indented JSON
func (l *ResourceList) ToJSONIndent() ([]byte, error) {
	data := make([]map[string]string, 0, len(l.items))
	for _, item := range l.items {
		data = append(data, item.Fields)
	}
	return json.MarshalIndent(data, "", "  ")
}
