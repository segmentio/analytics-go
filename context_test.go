package analytics

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMakeJsonContextNil(t *testing.T) {
	c := Context{}
	p := makeJsonContext(c)

	if p != nil {
		t.Error("zero-value context should be interpreted as a nil JSON object")
	}
}

func TestMakeJsonContextNonNil(t *testing.T) {
	c := Context{Locale: "en_US"}
	p := makeJsonContext(c)

	if p == nil {
		t.Error("non-zero-value context should not be interpreted as a nil JSON object")
	} else if !reflect.DeepEqual(*p, c) {
		t.Error("JSON object should equal the original context")
	}
}

func TestContextMarshalJSONLibrary(t *testing.T) {
	c := Context{
		Library: LibraryInfo{
			Name: "testing",
		},
	}

	if b, err := json.Marshal(c); err != nil {
		t.Error("marshaling context object failed:", err)

	} else if s := string(b); s != `{"library":{"name":"testing"}}` {
		t.Error("invalid marshaled representation of context:", s)
	}
}

func TestContextMarshalJSONExtra(t *testing.T) {
	c := Context{
		Extra: map[string]interface{}{
			"answer": 42,
		},
	}

	if b, err := json.Marshal(c); err != nil {
		t.Error("marshaling context object failed:", err)

	} else if s := string(b); s != `{"answer":42}` {
		t.Error("invalid marshaled representation of context:", s)
	}
}
