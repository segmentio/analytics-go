package analytics

import (
	"encoding/json"
	"testing"
)

func TestAppInfoBuildInt(t *testing.T) {
	data := []byte(`{"build": 698, "name": "ACME", "namespace": "com.acme.app", "version": "2.7.1"}`)

	var info AppInfo

	if err := json.Unmarshal(data, &info); err != nil {
		t.Error("unmarshalling app info json failed:", err)
	}

	if string(info.Build) != "698" {
		t.Errorf("Not equal: \nexpected: %s\nactual: %s", "698", string(info.Build))
	}
}

func TestAppInfoBuildString(t *testing.T) {
	data := []byte(`{"build": "698", "name": "ACME", "namespace": "com.acme.app", "version": "2.7.1"}`)

	var info AppInfo

	if err := json.Unmarshal(data, &info); err != nil {
		t.Error("unmarshalling app info json failed:", err)
	}

	if string(info.Build) != "698" {
		t.Errorf("Not equal: \nexpected: %s\nactual: %s", "698", string(info.Build))
	}
}

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

func TestContextMarshalJSONExtra(t *testing.T) {
	c := Context{
		Extra: map[string]interface{}{
			"answer": 42,
		},
	}

	if b, err := json.Marshal(c); err != nil {
		t.Error("marshalling context object failed:", err)

	} else if s := string(b); s != `{"answer":42}` {
		t.Error("invalid marshaled representation of context:", s)
	}
}
