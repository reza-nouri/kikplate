package plategenerator

import (
	"fmt"
	"strconv"
	"strings"
)

func CoerceByType(key, typ string, val any) (any, error) {
	typ = strings.ToLower(strings.TrimSpace(typ))
	switch typ {
	case "", "string", "enum":
		return fmt.Sprintf("%v", val), nil
	case "bool", "boolean":
		if b, ok := AsBool(val); ok {
			return b, nil
		}
		return nil, fmt.Errorf("%q must be a boolean", key)
	case "int", "integer":
		i, ok := AsInt64(val)
		if !ok {
			return nil, fmt.Errorf("%q must be an integer", key)
		}
		return i, nil
	case "number", "float", "float64":
		f, ok := AsFloat64(val)
		if !ok {
			return nil, fmt.Errorf("%q must be a number", key)
		}
		return f, nil
	default:
		return val, nil
	}
}

func AsBool(v any) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		s := strings.TrimSpace(strings.ToLower(t))
		switch s {
		case "true", "1", "yes", "y", "on":
			return true, true
		case "false", "0", "no", "n", "off":
			return false, true
		}
	case int:
		return t != 0, true
	case int64:
		return t != 0, true
	case float64:
		return t != 0, true
	}
	return false, false
}

func AsInt64(v any) (int64, bool) {
	switch t := v.(type) {
	case int:
		return int64(t), true
	case int64:
		return t, true
	case float64:
		if float64(int64(t)) == t {
			return int64(t), true
		}
		return 0, false
	case string:
		i, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		if err != nil {
			return 0, false
		}
		return i, true
	}
	return 0, false
}

func AsFloat64(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}

func IsEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}
