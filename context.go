package analytics

import (
	"encoding/json"
	"net"
	"reflect"
	"strings"
)

// This type provides the representation of the `context` object as defined in
// https://segment.com/docs/spec/common/#context
type Context struct {
	App       AppInfo                `json:"app"`
	Campaign  CampaignInfo           `json:"campaign"`
	Device    DeviceInfo             `json:"device"`
	IP        net.IP                 `json:"ip,omitempty"`
	Library   LibraryInfo            `json:"library"`
	Location  LocationInfo           `json:"location"`
	Network   NetworkInfo            `json:"network"`
	OS        OSInfo                 `json:"os"`
	Page      PageInfo               `json:"page"`
	Referrer  ReferrerInfo           `json:"referrer"`
	Screen    ScreenInfo             `json:"screen"`
	Locale    string                 `json:"locale,omitempty"`
	Timezone  string                 `json:"timezone,omitempty"`
	UserAgent string                 `json:"userAgent,omitempty"`
	Traits    map[string]interface{} `json:"traits"`

	// This map is used to allow extensions to the context specifications that
	// may not be documented or could be introduced in the future.
	// The fields of this map are inlined in the serialized context object,
	// there is no actual "extra" field in the JSON representation.
	Extra map[string]interface{} `json:"-"`
}

// This type provides the representation of the `context.app` object as defined
// in https://segment.com/docs/spec/common/#context
type AppInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Build   string `json:"build,omitempty"`
}

// This type provides the representation of the `context.campaign` object as
// defined in https://segment.com/docs/spec/common/#context
type CampaignInfo struct {
	Name    string `json:"name,omitempty"`
	Source  string `json:"source,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Term    string `json:"term,omitempty"`
	Content string `json:"content,omitempty"`
}

// This type provides the representation of the `context.device` object as
// defined in https://segment.com/docs/spec/common/#context
type DeviceInfo struct {
	Id           string `json:"id,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	Version      string `json:"version,omitempty"`
}

// This type provides the representation of the `context.library` object as
// defined in https://segment.com/docs/spec/common/#context
type LibraryInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// This type provides the representation of the `context.location` object as
// defined in https://segment.com/docs/spec/common/#context
type LocationInfo struct {
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Speed     float64 `json:"speed,omitempty"`
}

// This type provides the representation of the `context.network` object as
// defined in https://segment.com/docs/spec/common/#context
type NetworkInfo struct {
	Bluetooth bool   `json:"bluetooth,omitempty"`
	Cellular  bool   `json:"cellular,omitempty"`
	WIFI      bool   `json:"wifi,omitempty"`
	Carrier   string `json:"carrier,omitempty"`
}

// This type provides the representation of the `context.os` object as defined
// in https://segment.com/docs/spec/common/#context
type OSInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// This type provides the representation of the `context.page` object as
// defined in https://segment.com/docs/spec/common/#context
type PageInfo struct {
	Hash     string `json:"hash,omitempty"`
	Path     string `json:"path,omitempty"`
	Referrer string `json:"referrer,omitempty"`
	Search   string `json:"search,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
}

// This type provides the representation of the `context.referrer` object as
// defined in https://segment.com/docs/spec/common/#context
type ReferrerInfo struct {
	Type string `json:"json,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Link string `json:"link,omitempty"`
}

// This type provides the representation of the `context.screen` object as
// defined in https://segment.com/docs/spec/common/#context
type ScreenInfo struct {
	Density int `json:"density,omitempty"`
	Width   int `json:"width,omitempty"`
	Height  int `json:"height,omitempty"`
}

// Satisfy the `json.Marshaler` interface. We have to flatten out the `Extra`
// field but the standard json package doesn't support it yet.
// Implementing this interface allows us to override the default marshaling of
// the context object and to the inlining ourselves.
//
// Related discussion: https://github.com/golang/go/issues/6213
func (ctx Context) MarshalJSON() ([]byte, error) {
	t := reflect.TypeOf(ctx)
	v := reflect.ValueOf(ctx)
	n := v.NumField()
	m := make(map[string]interface{}, n+len(ctx.Extra))

	// Copy the `Extra` map into the map representation of the context, it is
	// important to do this operation before going through the actual struct
	// fields so the latter take precendence and override duplicated values
	// that would be set in the extensions.
	for name, value := range ctx.Extra {
		m[name] = value
	}

	// Imitate what what the JSON package would do when serializing a struct
	// value.
	for i := 0; i != n; i++ {
		name, omitempty := parseJsonTag(t.Field(i).Tag.Get("json"))
		value := v.Field(i)

		if name != "-" && !(omitempty && isEmptyValue(value)) {
			m[name] = value.Interface()
		}
	}

	return json.Marshal(m)
}

// Parses a JSON tag the way the json package would do it, returing the expected
// name of the field once serialized and if empty values should be omitted.
func parseJsonTag(tag string) (name string, omitempty bool) {
	args := strings.Split(tag, ",")

	if len(args) == 0 {
		name = f.Name
	} else {
		name = args[0]
	}

	if len(args) > 1 && args[1] == "omitempty" {
		omitempty = true
	}

	return
}

// This function was copied from https://golang.org/src/encoding/json/encode.go#L282
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
