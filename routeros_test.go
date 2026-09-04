package routeros

import (
	"testing"
	"time"
)

func TestEncodeLength(t *testing.T) {
	tests := []struct {
		length   int
		expected []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x80}},
		{255, []byte{0x80, 0xFF}},
		{256, []byte{0x81, 0x00}},
		{16383, []byte{0xBF, 0xFF}},
		{16384, []byte{0xC0, 0x40, 0x00}},
		{2097151, []byte{0xDF, 0xFF, 0xFF}},
		{2097152, []byte{0xE0, 0x20, 0x00, 0x00}},
	}

	for _, test := range tests {
		result, err := EncodeLength(test.length)
		if err != nil {
			t.Errorf("EncodeLength(%d) error: %v", test.length, err)
			continue
		}
		if !equalBytes(result, test.expected) {
			t.Errorf("EncodeLength(%d) = %v, expected %v", test.length, result, test.expected)
		}
	}
}

func TestDecodeLength(t *testing.T) {
	tests := []struct {
		bytes    []byte
		expected int
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x7F}, 127},
		{[]byte{0x80, 0x80}, 128},
		{[]byte{0x80, 0xFF}, 255},
		{[]byte{0x81, 0x00}, 256},
		{[]byte{0xBF, 0xFF}, 16383},
		{[]byte{0xC0, 0x40, 0x00}, 16384},
		{[]byte{0xDF, 0xFF, 0xFF}, 2097151},
		{[]byte{0xE0, 0x20, 0x00, 0x00}, 2097152},
	}

	for _, test := range tests {
		result, err := DecodeLength(&mockReader{data: test.bytes})
		if err != nil {
			t.Errorf("DecodeLength(%v) error: %v", test.bytes, err)
			continue
		}
		if result != test.expected {
			t.Errorf("DecodeLength(%v) = %d, expected %d", test.bytes, result, test.expected)
		}
	}
}

type mockReader struct {
	data []byte
	pos  int
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.data) {
		return 0, nil
	}
	n = copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseWord(t *testing.T) {
	tests := []struct {
		word     string
		expected Word
	}{
		{"/interface/print", Word{Type: WordTypeCommand, Value: "/interface/print"}},
		{"=name=ether1", Word{Type: WordTypeAPIAttribute, Key: "name", Value: "ether1"}},
		{"=disabled", Word{Type: WordTypeAPIAttribute, Key: "disabled", Value: ""}},
		{".tag=1", Word{Type: WordTypeAPIAttribute, Value: ".tag=1"}},
		{"?name=ether1", Word{Type: WordTypeQuery, Value: "?name=ether1"}},
		{"name=value", Word{Type: WordTypeAttribute, Key: "name", Value: "value"}},
	}

	for _, test := range tests {
		result := ParseWord(test.word)
		if result.Type != test.expected.Type {
			t.Errorf("ParseWord(%q).Type = %v, expected %v", test.word, result.Type, test.expected.Type)
		}
		if result.Key != test.expected.Key {
			t.Errorf("ParseWord(%q).Key = %v, expected %v", test.word, result.Key, test.expected.Key)
		}
		if result.Value != test.expected.Value {
			t.Errorf("ParseWord(%q).Value = %v, expected %v", test.word, result.Value, test.expected.Value)
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"1s", time.Second},
		{"1m", time.Minute},
		{"1h", time.Hour},
		{"1d", 24 * time.Hour},
		{"1d2h3m4s", 24*time.Hour + 2*time.Hour + 3*time.Minute + 4*time.Second},
		{"100ms", 100 * time.Millisecond},
	}

	for _, test := range tests {
		result, err := ParseDuration(test.input)
		if err != nil {
			t.Errorf("ParseDuration(%q) error: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("ParseDuration(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"yes", true},
		{"no", false},
		{"1", true},
		{"0", false},
		{"TRUE", true},
		{"FALSE", false},
	}

	for _, test := range tests {
		result := ParseBool(test.input)
		if result != test.expected {
			t.Errorf("ParseBool(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestQueryBuilder(t *testing.T) {
	// Test simple query
	q := NewQuery("/interface/print").
		Equal("name", "ether1").
		PropList("name", "type")

	words := q.Build()
	if len(words) == 0 {
		t.Error("Query build returned empty words")
	}

	if words[0] != "/interface/print" {
		t.Errorf("First word should be command, got %s", words[0])
	}
}

func TestQueryBuilderConditions(t *testing.T) {
	q := NewQuery("/ip/firewall/filter/print").
		Equal("chain", "input").
		Equal("action", "accept").
		And()

	words := q.Build()
	if len(words) < 4 {
		t.Errorf("Expected at least 4 words, got %d", len(words))
	}
}

func TestConfigValidation(t *testing.T) {
	// Test missing address
	config := NewConfig().WithCredentials("admin", "password")
	if err := config.Validate(); err == nil {
		t.Error("Expected error for missing address")
	}

	// Test missing username
	config = NewConfig().WithAddress("192.168.88.1")
	if err := config.Validate(); err == nil {
		t.Error("Expected error for missing username")
	}

	// Test valid config
	config = NewConfig().
		WithAddress("192.168.88.1").
		WithCredentials("admin", "password")
	if err := config.Validate(); err != nil {
		t.Errorf("Unexpected error for valid config: %v", err)
	}
}

func TestConfigMethods(t *testing.T) {
	config := NewConfig().
		WithAddress("192.168.88.1").
		WithPort(8729).
		WithCredentials("admin", "password").
		WithTLS().
		WithTimeout(60 * time.Second).
		WithAutoReconnect(true).
		WithSkipVerify(true)

	if config.Address != "192.168.88.1" {
		t.Errorf("Address = %s, expected 192.168.88.1", config.Address)
	}
	if config.Port != 8729 {
		t.Errorf("Port = %d, expected 8729", config.Port)
	}
	if !config.TLS {
		t.Error("TLS should be true")
	}
	if !config.AutoReconnect {
		t.Error("AutoReconnect should be true")
	}
	if !config.SkipVerify {
		t.Error("SkipVerify should be true")
	}
}

func TestErrors(t *testing.T) {
	// Test RouterError
	routerErr := NewRouterError(map[string]string{
		"message":  "test error",
		"category": "2",
		"tag":      "1",
	})

	if routerErr.Message != "test error" {
		t.Errorf("RouterError.Message = %s, expected 'test error'", routerErr.Message)
	}
	if routerErr.Category != 2 {
		t.Errorf("RouterError.Category = %d, expected 2", routerErr.Category)
	}

	// Test ConnectionError
	connErr := NewConnectionError("192.168.88.1", 8728, ErrNotConnected)
	if connErr.Addr != "192.168.88.1" {
		t.Errorf("ConnectionError.Addr = %s, expected 192.168.88.1", connErr.Addr)
	}
}

func TestResponseParser(t *testing.T) {
	responses := []*Response{
		{Type: ResponseReply, Data: map[string]string{"name": "ether1", "type": "ethernet"}},
		{Type: ResponseReply, Data: map[string]string{"name": "ether2", "type": "ethernet"}},
		{Type: ResponseDone, Done: true},
	}

	parser := NewResponseParser(responses)

	if parser.Count() != 2 {
		t.Errorf("Count = %d, expected 2", parser.Count())
	}

	if !parser.Done() {
		t.Error("Done should be true")
	}

	if parser.Value("name") != "ether1" {
		t.Errorf("Value('name') = %s, expected 'ether1'", parser.Value("name"))
	}
}

func TestResource(t *testing.T) {
	data := map[string]string{
		".id":       "*1",
		"name":      "ether1",
		"type":      "ethernet",
		"disabled":  "false",
		"mtu":       "1500",
	}

	resource := NewResource(data)

	if resource.ID != "*1" {
		t.Errorf("ID = %s, expected '*1'", resource.ID)
	}
	if resource.Name != "ether1" {
		t.Errorf("Name = %s, expected 'ether1'", resource.Name)
	}
	if resource.Get("type") != "ethernet" {
		t.Errorf("Get('type') = %s, expected 'ethernet'", resource.Get("type"))
	}

	mtu, err := resource.GetInt("mtu")
	if err != nil {
		t.Errorf("GetInt error: %v", err)
	}
	if mtu != 1500 {
		t.Errorf("GetInt('mtu') = %d, expected 1500", mtu)
	}

	if !resource.GetBool("disabled") {
		t.Error("GetBool('disabled') should be true for 'false' string")
	}
}

func TestResourceList(t *testing.T) {
	responses := []*Response{
		{Type: ResponseReply, Data: map[string]string{".id": "*1", "name": "ether1"}},
		{Type: ResponseReply, Data: map[string]string{".id": "*2", "name": "ether2"}},
		{Type: ResponseDone, Done: true},
	}

	list := NewResourceList(responses)

	if list.Len() != 2 {
		t.Errorf("Len = %d, expected 2", list.Len())
	}

	first := list.First()
	if first.Name != "ether1" {
		t.Errorf("First().Name = %s, expected 'ether1'", first.Name)
	}

	ids := list.IDs()
	if len(ids) != 2 {
		t.Errorf("IDs length = %d, expected 2", len(ids))
	}

	names := list.Names()
	if len(names) != 2 {
		t.Errorf("Names length = %d, expected 2", len(names))
	}

	// Find by ID
	found := list.Get("*2")
	if found == nil {
		t.Error("Get('*2') should not be nil")
	}

	// Find by name
	found = list.FindByName("ether2")
	if found == nil {
		t.Error("FindByName('ether2') should not be nil")
	}
}

func TestPathBuilder(t *testing.T) {
	path := IP().Address()
	if path.String() != "/ip/address" {
		t.Errorf("Path = %s, expected '/ip/address'", path.String())
	}

	path = Firewall().Filter()
	if path.String() != "/ip/firewall/filter" {
		t.Errorf("Path = %s, expected '/ip/firewall/filter'", path.String())
	}
}
