package routeros

import "testing"

func TestParseAPITagAttribute(t *testing.T) {
	word := ParseWord(".tag=42")
	if word.Type != WordTypeAPIAttribute {
		t.Fatalf("type = %v, want WordTypeAPIAttribute", word.Type)
	}
	if word.Key != "tag" || word.Value != "42" {
		t.Fatalf("parsed API tag = key %q value %q, want tag/42", word.Key, word.Value)
	}

	response := ParseSentence(&Sentence{Words: []Word{
		ParseWord("!re"),
		ParseWord(".tag=42"),
		ParseWord("name=ether1"),
	}})
	if response.Tag != "42" {
		t.Fatalf("response tag = %q, want 42", response.Tag)
	}
	if response.Data["name"] != "ether1" {
		t.Fatalf("response data = %#v, want name=ether1", response.Data)
	}
}
