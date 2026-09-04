package routeros

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Protocol constants
const (
	// Control bytes for the RouterOS API protocol
	WordLengthByte = 0x00 // Word length marker

	// Maximum word length (268435455 bytes = 0x0FFFFFFF)
	MaxWordLength = 0x0FFFFFFF

	// Sentence end marker
	SentenceEnd = ""
)

// WordType represents the type of a word in the API protocol
type WordType int

const (
	WordTypeCommand    WordType = iota // Command word (starts with /)
	WordTypeAttribute                  // Attribute word (name=value)
	WordTypeAPIAttribute               // API attribute (starts with =, ., or <)
	WordTypeQuery                      // Query word
	WordTypeEnd                        // End of sentence marker
)

// Word represents a single word in the API protocol
type Word struct {
	Type  WordType
	Key   string // For attributes: the key name
	Value string // For attributes: the value; for commands: the command string
}

// Sentence represents a complete sentence in the API protocol
type Sentence struct {
	Words []Word
	Tag   string // Optional tag for tracking responses
}

// ResponseType represents the type of response
type ResponseType string

const (
	ResponseReply   ResponseType = "!re"  // Reply with data
	ResponseDone    ResponseType = "!done" // Command completed
	ResponseTrap    ResponseType = "!trap" // Error/trap
	ResponseFatal   ResponseType = "!fatal" // Fatal error
	ResponseData    ResponseType = ""    // Data line (no prefix)
)

// Response represents a parsed response from the router
type Response struct {
	Type    ResponseType
	Tag     string
	Data    map[string]string
	Message string // For error messages
	Done    bool   // Indicates if this is the final response
}

// Protocol handles low-level protocol operations
type Protocol struct {
	reader io.Reader
	writer io.Writer
}

// NewProtocol creates a new Protocol instance
func NewProtocol(reader io.Reader, writer io.Writer) *Protocol {
	return &Protocol{
		reader: reader,
		writer: writer,
	}
}

// EncodeLength encodes a length value according to the RouterOS API protocol
func EncodeLength(length int) ([]byte, error) {
	if length < 0 || length > MaxWordLength {
		return nil, fmt.Errorf("invalid word length: %d", length)
	}

	if length < 0x80 {
		return []byte{byte(length)}, nil
	}

	if length < 0x4000 {
		return []byte{
			byte((length >> 8) | 0x80),
			byte(length & 0xFF),
		}, nil
	}

	if length < 0x200000 {
		return []byte{
			byte((length >> 16) | 0xC0),
			byte((length >> 8) & 0xFF),
			byte(length & 0xFF),
		}, nil
	}

	return []byte{
		byte((length >> 24) | 0xE0),
		byte((length >> 16) & 0xFF),
		byte((length >> 8) & 0xFF),
		byte(length & 0xFF),
	}, nil
}

// DecodeLength decodes a length value according to the RouterOS API protocol
func DecodeLength(reader io.Reader) (int, error) {
	var firstByte [1]byte
	_, err := io.ReadFull(reader, firstByte[:])
	if err != nil {
		return 0, err
	}

	b := firstByte[0]

	switch {
	case b&0xE0 == 0xE0:
		// 4-byte length (most significant 4 bits of first byte + 3 more bytes)
		var remaining [3]byte
		_, err = io.ReadFull(reader, remaining[:])
		if err != nil {
			return 0, err
		}
		length := int(b&0x1F) << 24
		length |= int(remaining[0]) << 16
		length |= int(remaining[1]) << 8
		length |= int(remaining[2])
		return length, nil

	case b&0xC0 == 0xC0:
		// 3-byte length
		var remaining [2]byte
		_, err = io.ReadFull(reader, remaining[:])
		if err != nil {
			return 0, err
		}
		length := int(b&0x3F) << 16
		length |= int(remaining[0]) << 8
		length |= int(remaining[1])
		return length, nil

	case b&0x80 == 0x80:
		// 2-byte length
		var secondByte [1]byte
		_, err = io.ReadFull(reader, secondByte[:])
		if err != nil {
			return 0, err
		}
		length := int(b&0x7F) << 8
		length |= int(secondByte[0])
		return length, nil

	default:
		// Single-byte length
		return int(b), nil
	}
}

// WriteWord writes a single word to the connection
func (p *Protocol) WriteWord(word string) error {
	// Encode and write length
	lengthBytes, err := EncodeLength(len(word))
	if err != nil {
		return err
	}

	_, err = p.writer.Write(lengthBytes)
	if err != nil {
		return err
	}

	// Write word content
	if len(word) > 0 {
		_, err = p.writer.Write([]byte(word))
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteSentence writes a complete sentence to the connection
func (p *Protocol) WriteSentence(words ...string) error {
	for _, word := range words {
		if err := p.WriteWord(word); err != nil {
			return err
		}
	}
	// Write sentence end marker (zero-length word)
	return p.WriteWord("")
}

// ReadWord reads a single word from the connection
func (p *Protocol) ReadWord() (string, error) {
	length, err := DecodeLength(p.reader)
	if err != nil {
		return "", err
	}

	if length == 0 {
		return "", nil // End of sentence
	}

	if length > MaxWordLength {
		return "", fmt.Errorf("word length exceeds maximum: %d", length)
	}

	buffer := make([]byte, length)
	_, err = io.ReadFull(p.reader, buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

// ReadSentence reads a complete sentence from the connection
func (p *Protocol) ReadSentence() (*Sentence, error) {
	sentence := &Sentence{
		Words: make([]Word, 0),
	}

	for {
		word, err := p.ReadWord()
		if err != nil {
			return nil, err
		}

		if word == "" {
			// End of sentence
			break
		}

		parsedWord := ParseWord(word)
		sentence.Words = append(sentence.Words, parsedWord)
	}

	return sentence, nil
}

// ParseWord parses a word string into a Word struct
func ParseWord(word string) Word {
	w := Word{Value: word}

	switch {
	case len(word) == 0:
		w.Type = WordTypeEnd

	case word[0] == '/':
		w.Type = WordTypeCommand
		w.Value = word

	case word[0] == '=':
		w.Type = WordTypeAPIAttribute
		if idx := bytes.IndexByte([]byte(word[1:]), '='); idx >= 0 {
			// Format: =key=value
			w.Key = word[1 : idx+1]
			w.Value = word[idx+2:]
		} else {
			// Format: =key (no value)
			w.Key = word[1:]
			w.Value = ""
		}

	case word[0] == '.':
		w.Type = WordTypeAPIAttribute
		w.Value = word

	case word[0] == '<' || word[0] == '>' || word[0] == '-':
		w.Type = WordTypeQuery
		w.Value = word

	default:
		// Check for key=value format
		if idx := bytes.IndexByte([]byte(word), '='); idx >= 0 {
			w.Type = WordTypeAttribute
			w.Key = word[:idx]
			w.Value = word[idx+1:]
		} else {
			w.Type = WordTypeCommand
		}
	}

	return w
}

// ParseSentence converts a Sentence to a Response
func ParseSentence(sentence *Sentence) *Response {
	response := &Response{
		Data: make(map[string]string),
	}

	for _, word := range sentence.Words {
		switch word.Type {
		case WordTypeCommand:
			switch word.Value {
			case "!re":
				response.Type = ResponseReply
			case "!done":
				response.Type = ResponseDone
				response.Done = true
			case "!trap":
				response.Type = ResponseTrap
			case "!fatal":
				response.Type = ResponseFatal
			default:
				response.Data[word.Value] = ""
			}

		case WordTypeAttribute, WordTypeAPIAttribute:
			if word.Key == "tag" {
				response.Tag = word.Value
			} else if word.Key == "message" {
				response.Message = word.Value
			} else if word.Key != "" {
				response.Data[word.Key] = word.Value
			}
		}
	}

	return response
}

// IsResponseType checks if a word is a response type marker
func IsResponseType(word string) bool {
	return word == "!re" || word == "!done" || word == "!trap" || word == "!fatal"
}

// EncodeWord encodes a key-value pair into an API attribute word
func EncodeWord(key, value string) string {
	return "=" + key + "=" + value
}

// EncodeQuery encodes a query condition
func EncodeQuery(op string, key string, value string) string {
	switch op {
	case "=", "==":
		return "?=" + key + "=" + value
	case "!=", "<>":
		return "?!=" + key + "=" + value
	case "<":
		return "?<" + key + "=" + value
	case ">":
		return "?<" + key + "=" + value
	case "<=":
		return "?<=" + key + "=" + value
	case ">=":
		return "?>=" + key + "=" + value
	default:
		return "?" + key + "=" + value
	}
}

// EncodeNumber encodes a number for RouterOS API
func EncodeNumber(num int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}

// DecodeNumber decodes a number from RouterOS API
func DecodeNumber(data []byte) int64 {
	var num int64
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.BigEndian, &num)
	return num
}
