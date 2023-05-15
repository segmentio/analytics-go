package journify

import (
	"encoding/json"
	"testing"
)

func TestContextMarshalJSONLibrary(t *testing.T) {
	c := Context{
		Library: LibraryInfo{
			Name: "testing",
		},
	}

	if b, err := json.Marshal(c); err != nil {
		t.Error("marshalling context object failed:", err)

	} else if s := string(b); s != `{"library":{"name":"testing"}}` {
		t.Error("invalid marshaled representation of context:", s)
	}
}
