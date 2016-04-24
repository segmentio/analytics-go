package analytics

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestMessageIdDefault(t *testing.T) {
	if id := makeMessageId("", "42"); id != "42" {
		t.Error("invalid default message id:", id)
	}
}

func TestMessageIdNonDefault(t *testing.T) {
	if id := makeMessageId("A", "42"); id != "A" {
		t.Error("invalid non-default message id:", id)
	}
}

func validateSerizable(t string, m Message) (v interface{}, err error) {
	v = m.serializable("", time.Now())

	v0 := reflect.ValueOf(m)
	v1 := reflect.ValueOf(v)

	t0 := v0.Type()
	t1 := v1.Type()

	for i := 0; i != v0.NumField(); i++ {
		ft0 := t0.Field(i)
		ft1, ok := t1.FieldByName(ft0.Name)

		if !ok {
			err = fmt.Errorf("missing '%s' field in serializable value:", ft0.Name)
			return
		}

		fv0 := v0.Field(i)
		fv1 := v1.FieldByIndex(ft1.Index)

		if len(ft1.Tag.Get("json")) == 0 {
			err = fmt.Errorf("missing `json` tag in serializable field '%s'", ft1.Name)
			return
		}

		x0 := fv0.Interface()
		x1 := fv1.Interface()

		switch x := x0.(type) {
		case time.Time:
			// Special case for timestamps because they get converted to strings.
			x0 = formatTime(x)
		case Context:
			// Special case for contexts as they are converted to pointers when
			// being serialized.
			x0 = makeJsonContext(x)
		}

		if !reflect.DeepEqual(x0, x1) {
			err = fmt.Errorf("invalid field value for '%s': expected %#v but found %#v", ft1.Name, x0, x1)
			return
		}
	}

	if f := v1.FieldByName("Type").Interface(); !reflect.DeepEqual(f, t) {
		err = fmt.Errorf("invalid field value for 'Type': %#v", f)
		return
	}

	if n0, n1 := (v0.NumField() + 1), v1.NumField(); n0 != n1 {
		err = fmt.Errorf("invalid field count in serializable value: expected %d but got %d", n0, n1)
		return
	}

	return
}
