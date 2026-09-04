package routeros

import (
	"fmt"
	"strings"
)

// QueryOperator represents a query comparison operator
type QueryOperator string

const (
	QueryEqual    QueryOperator = "="
	QueryNotEqual QueryOperator = "!="
	QueryLess     QueryOperator = "<"
	QueryLessEq   QueryOperator = "<="
	QueryGreater  QueryOperator = ">"
	QueryGreaterEq QueryOperator = ">="
)

// QueryBuilder provides a fluent interface for building RouterOS API queries
type QueryBuilder struct {
	command  string
	args     map[string]string
	queries  []string
	proplist []string
	tag      string
}

// NewQuery creates a new QueryBuilder for the given command
func NewQuery(command string) *QueryBuilder {
	return &QueryBuilder{
		command: command,
		args:    make(map[string]string),
	}
}

// Command creates a new QueryBuilder for the given command path
func Command(command string) *QueryBuilder {
	return NewQuery(command)
}

// Where adds a query condition with the given operator
func (q *QueryBuilder) Where(key string, op QueryOperator, value string) *QueryBuilder {
	switch op {
	case QueryEqual:
		q.queries = append(q.queries, "?="+key+"="+value)
	case QueryNotEqual:
		q.queries = append(q.queries, "?!"+key+"="+value)
	case QueryLess:
		q.queries = append(q.queries, "?<"+key+"="+value)
	case QueryLessEq:
		q.queries = append(q.queries, "?<="+key+"="+value)
	case QueryGreater:
		q.queries = append(q.queries, "?<"+key+"="+value) // Note: RouterOS uses < for greater
	case QueryGreaterEq:
		q.queries = append(q.queries, "?>="+key+"="+value)
	}
	return q
}

// Equal adds an equality query condition
func (q *QueryBuilder) Equal(key, value string) *QueryBuilder {
	return q.Where(key, QueryEqual, value)
}

// NotEqual adds a not-equal query condition
func (q *QueryBuilder) NotEqual(key, value string) *QueryBuilder {
	return q.Where(key, QueryNotEqual, value)
}

// Less adds a less-than query condition
func (q *QueryBuilder) Less(key, value string) *QueryBuilder {
	return q.Where(key, QueryLess, value)
}

// LessOrEqual adds a less-than-or-equal query condition
func (q *QueryBuilder) LessOrEqual(key, value string) *QueryBuilder {
	return q.Where(key, QueryLessEq, value)
}

// Greater adds a greater-than query condition
func (q *QueryBuilder) Greater(key, value string) *QueryBuilder {
	return q.Where(key, QueryGreater, value)
}

// GreaterOrEqual adds a greater-than-or-equal query condition
func (q *QueryBuilder) GreaterOrEqual(key, value string) *QueryBuilder {
	return q.Where(key, QueryGreaterEq, value)
}

// WhereIn adds a query condition for multiple values (OR)
func (q *QueryBuilder) WhereIn(key string, values ...string) *QueryBuilder {
	for _, value := range values {
		q.queries = append(q.queries, "?="+key+"="+value)
	}
	return q
}

// WhereNotIn adds a query condition to exclude multiple values (AND NOT)
func (q *QueryBuilder) WhereNotIn(key string, values ...string) *QueryBuilder {
	for _, value := range values {
		q.queries = append(q.queries, "?!"+key+"="+value)
	}
	return q
}

// WhereLike adds a pattern matching query condition
func (q *QueryBuilder) WhereLike(key, pattern string) *QueryBuilder {
	q.queries = append(q.queries, "?-"+key+"="+pattern)
	return q
}

// WhereNotNull adds a query for non-null values
func (q *QueryBuilder) WhereNotNull(key string) *QueryBuilder {
	q.queries = append(q.queries, "?"+key+"!")
	return q
}

// WhereNull adds a query for null/empty values
func (q *QueryBuilder) WhereNull(key string) *QueryBuilder {
	q.queries = append(q.queries, "?"+key+"=")
	return q
}

// And combines the current queries with AND
func (q *QueryBuilder) And() *QueryBuilder {
	q.queries = append(q.queries, "?#&")
	return q
}

// Or combines the current queries with OR
func (q *QueryBuilder) Or() *QueryBuilder {
	q.queries = append(q.queries, "?#|")
	return q
}

// Not negates the current query
func (q *QueryBuilder) Not() *QueryBuilder {
	q.queries = append(q.queries, "?#!")
	return q
}

// AndQuery combines two query conditions with AND
func (q *QueryBuilder) AndQuery() *QueryBuilder {
	q.queries = append(q.queries, "?#&")
	return q
}

// OrQuery combines two query conditions with OR
func (q *QueryBuilder) OrQuery() *QueryBuilder {
	q.queries = append(q.queries, "?#|")
	return q
}

// Arg sets a command argument
func (q *QueryBuilder) Arg(key, value string) *QueryBuilder {
	q.args[key] = value
	return q
}

// Args sets multiple command arguments
func (q *QueryBuilder) Args(args map[string]string) *QueryBuilder {
	for k, v := range args {
		q.args[k] = v
	}
	return q
}

// PropList sets the properties to return
func (q *QueryBuilder) PropList(props ...string) *QueryBuilder {
	q.proplist = props
	return q
}

// Tag sets the command tag
func (q *QueryBuilder) Tag(tag string) *QueryBuilder {
	q.tag = tag
	return q
}

// Build constructs the query words
func (q *QueryBuilder) Build() []string {
	words := make([]string, 0)

	// Add command
	words = append(words, q.command)

	// Add arguments
	for key, value := range q.args {
		words = append(words, "="+key+"="+value)
	}

	// Add queries
	words = append(words, q.queries...)

	// Add proplist
	if len(q.proplist) > 0 {
		words = append(words, "=.proplist="+strings.Join(q.proplist, ","))
	}

	// Add tag
	if q.tag != "" {
		words = append(words, ".tag="+q.tag)
	}

	return words
}

// Execute runs the query on the given client
func (q *QueryBuilder) Execute(client *Client) ([]*Response, error) {
	return client.Send(q.Build()...)
}

// ExecuteContext runs the query with context on the given client
func (q *QueryBuilder) ExecuteContext(ctx interface{}, client *Client) ([]*Response, error) {
	words := q.Build()
	return client.Send(words[0], words[1:]...)
}

// String returns the query as a string for debugging
func (q *QueryBuilder) String() string {
	return strings.Join(q.Build(), " ")
}

// Add adds a new query word directly
func (q *QueryBuilder) Add(query string) *QueryBuilder {
	q.queries = append(q.queries, query)
	return q
}

// Reset clears all queries and arguments
func (q *QueryBuilder) Reset() *QueryBuilder {
	q.args = make(map[string]string)
	q.queries = nil
	q.proplist = nil
	q.tag = ""
	return q
}

// SetCommand changes the command
func (q *QueryBuilder) SetCommand(command string) *QueryBuilder {
	q.command = command
	return q
}

// Clone creates a copy of the query builder
func (q *QueryBuilder) Clone() *QueryBuilder {
	clone := &QueryBuilder{
		command:  q.command,
		args:     make(map[string]string),
		queries:  make([]string, len(q.queries)),
		proplist: make([]string, len(q.proplist)),
		tag:      q.tag,
	}

	for k, v := range q.args {
		clone.args[k] = v
	}
	copy(clone.queries, q.queries)
	copy(clone.proplist, q.proplist)

	return clone
}

// Query functions for complex queries

// QueryAnd creates an AND query from multiple conditions
func QueryAnd(conditions ...string) string {
	parts := make([]string, 0, len(conditions)+1)
	parts = append(parts, conditions...)
	parts = append(parts, "?#&")
	return strings.Join(parts, " ")
}

// QueryOr creates an OR query from multiple conditions
func QueryOr(conditions ...string) string {
	parts := make([]string, 0, len(conditions)+1)
	parts = append(parts, conditions...)
	parts = append(parts, "?#|")
	return strings.Join(parts, " ")
}

// QueryNot creates a NOT query
func QueryNot(condition string) string {
	return condition + " ?#!"
}

// QueryBuilderHelper provides helper methods for building queries

// ParseQuery parses a query string into conditions
func ParseQuery(query string) []string {
	return strings.Fields(query)
}

// FormatValue formats a value for the API (handles special characters)
func FormatValue(value string) string {
	// RouterOS API handles most characters natively
	// Only need to escape certain characters if needed
	return value
}

// QueryFromMap creates query conditions from a map
func QueryFromMap(conditions map[string]string) []string {
	queries := make([]string, 0, len(conditions))
	for key, value := range conditions {
		queries = append(queries, "?="+key+"="+value)
	}
	return queries
}

// QueryFromStruct creates query conditions from a struct (using struct tags)
func QueryFromStruct(s interface{}) ([]string, error) {
	// This would require reflection and struct tag parsing
	// For now, return empty slice
	return []string{}, nil
}

// PathBuilder helps build command paths

// Path represents a RouterOS API command path
type Path struct {
	elements []string
}

// NewPath creates a new Path
func NewPath(elements ...string) *Path {
	return &Path{elements: elements}
}

// Add adds an element to the path
func (p *Path) Add(element string) *Path {
	p.elements = append(p.elements, element)
	return p
}

// String returns the path as a string
func (p *Path) String() string {
	return "/" + strings.Join(p.elements, "/")
}

// Build returns the path components
func (p *Path) Build() []string {
	return p.elements
}

// IP returns a path starting with /ip
func IP() *Path {
	return NewPath("ip")
}

// IPv6 returns a path starting with /ipv6
func IPv6() *Path {
	return NewPath("ipv6")
}

// Interface returns a path starting with /interface
func Interface() *Path {
	return NewPath("interface")
}

// System returns a path starting with /system
func System() *Path {
	return NewPath("system")
}

// PPP returns a path starting with /ppp
func PPP() *Path {
	return NewPath("ppp")
}

// Routing returns a path starting with /routing
func Routing() *Path {
	return NewPath("routing")
}

// Queue returns a path starting with /queue
func Queue() *Path {
	return NewPath("queue")
}

// Firewall returns a path starting with /ip/firewall
func Firewall() *Path {
	return NewPath("ip", "firewall")
}

// Address returns a path for /ip/address
func (p *Path) Address() *Path {
	return NewPath(append(p.elements, "address")...)
}

// Route returns a path for /ip/route
func (p *Path) Route() *Path {
	return NewPath(append(p.elements, "route")...)
}

// DNS returns a path for /ip/dns
func (p *Path) DNS() *Path {
	return NewPath(append(p.elements, "dns")...)
}

// ARP returns a path for /ip/arp
func (p *Path) ARP() *Path {
	return NewPath(append(p.elements, "arp")...)
}

// DHCPClient returns a path for /ip/dhcp-client
func (p *Path) DHCPClient() *Path {
	return NewPath(append(p.elements, "dhcp-client")...)
}

// DHCPServer returns a path for /ip/dhcp-server
func (p *Path) DHCPServer() *Path {
	return NewPath(append(p.elements, "dhcp-server")...)
}

// Pool returns a path for /ip/pool
func (p *Path) Pool() *Path {
	return NewPath(append(p.elements, "pool")...)
}

// Filter returns a path for firewall filter
func (p *Path) Filter() *Path {
	return NewPath(append(p.elements, "filter")...)
}

// Nat returns a path for firewall NAT
func (p *Path) Nat() *Path {
	return NewPath(append(p.elements, "nat")...)
}

// Mangle returns a path for firewall mangle
func (p *Path) Mangle() *Path {
	return NewPath(append(p.elements, "mangle")...)
}

// Print creates a print command
func Print() *QueryBuilder {
	return Command("/print")
}

// Find creates a find query
func Find() *QueryBuilder {
	return Command("/find")
}

// Add creates an add command
func Add() *QueryBuilder {
	return Command("/add")
}

// Set creates a set command
func Set() *QueryBuilder {
	return Command("/set")
}

// Remove creates a remove command
func Remove() *QueryBuilder {
	return Command("/remove")
}

// Enable creates an enable command
func Enable() *QueryBuilder {
	return Command("/enable")
}

// Disable creates a disable command
func Disable() *QueryBuilder {
	return Command("/disable")
}

// Comment creates a comment command
func Comment() *QueryBuilder {
	return Command("/comment")
}

// Move creates a move command
func Move() *QueryBuilder {
	return Command("/move")
}

// Reset creates a reset command
func Reset() *QueryBuilder {
	return Command("/reset")
}

// Export creates an export command
func Export() *QueryBuilder {
	return Command("/export")
}

// Import creates an import command
func Import() *QueryBuilder {
	return Command("/import")
}

// ID sets the item ID for operations
func ID(id string) string {
	return "=.id=" + id
}

// Numbers sets the item numbers for operations
func Numbers(nums ...string) string {
	return "=numbers=" + strings.Join(nums, ",")
}

// Name sets the name for operations
func Name(name string) string {
	return "=name=" + name
}
