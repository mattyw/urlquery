package urlquery

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Marshal(s interface{}) (string, error) {
	values := url.Values{}
	st := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		fv := v.Field(i)
		value, err := marshalValue(fv.Interface())
		if err != nil {
			return "", err
		}
		key := f.Tag.Get("urlquery")
		if key == "" {
			key = strings.ToLower(f.Name)
		}
		values.Add(key, value)
	}
	return values.Encode(), nil
}

func marshalValue(i interface{}) (string, error) {
	switch t := i.(type) {
	case bool:
		return fmt.Sprintf("%t", t), nil
	case int:
		return fmt.Sprintf("%d", i), nil
	case string:
		return fmt.Sprintf("%s", t), nil
	case time.Time:
		return t.Format(time.RFC3339), nil
	default:
		return "", fmt.Errorf("unsupported type %v", t)
	}
}

func Unmarshal(in string, out interface{}) error {
	values, err := url.ParseQuery(in)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(out)
	switch v.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Map:
		unMarhshalTo(values, v.Elem())
	case reflect.Struct:
		return fmt.Errorf("can't deal with struct values. Use a pointer.")
	default:
		return fmt.Errorf("Unmarshal needs a map or a pointer to a struct.")
	}
	return nil
}

func unMarhshalTo(values url.Values, out reflect.Value) error {
	var empty reflect.Value
	for k, v := range values {
		fv := out.FieldByName(strings.Title(k))
		st := reflect.TypeOf(out)
		for i := 0; i < st.NumField(); i++ {
			f := st.Field(i)
			if f.Tag.Get("urlquery") != "" {
				k = f.Tag.Get("urlquery")
			}
		}
		if fv == empty {
			continue
		}
		switch fv.Interface().(type) {
		case bool:
			if v[0] == "true" {
				fv.SetBool(true)
			} else {
				fv.SetBool(false)
			}
		case int:
			in, err := strconv.Atoi(v[0])
			if err != nil {
				return err
			}
			fv.SetInt(int64(in))
		case string:
			fv.SetString(v[0])
		case time.Time:
			t, err := time.Parse(time.RFC3339, v[0])
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(t))
		}
	}
	return nil
}
