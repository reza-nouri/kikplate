package plategenerator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func EvalCondition(condition string, data map[string]any) bool {
	condition = strings.TrimSpace(condition)
	if condition == "" {
		return true
	}
	if result, ok := evalBoolExpr(condition, data); ok {
		return result
	}
	v, _ := resolveExprValue(condition, data)
	return truthy(v)
}

func evalBoolExpr(expr string, data map[string]any) (bool, bool) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false, false
	}

	expr = trimBalancedParens(expr)

	if parts := splitTopLevel(expr, "||"); len(parts) > 1 {
		for _, p := range parts {
			v, ok := evalBoolExpr(p, data)
			if !ok {
				return false, false
			}
			if v {
				return true, true
			}
		}
		return false, true
	}

	if parts := splitTopLevel(expr, "&&"); len(parts) > 1 {
		for _, p := range parts {
			v, ok := evalBoolExpr(p, data)
			if !ok {
				return false, false
			}
			if !v {
				return false, true
			}
		}
		return true, true
	}

	if strings.HasPrefix(expr, "!") {
		v, ok := evalBoolExpr(strings.TrimSpace(expr[1:]), data)
		if !ok {
			return false, false
		}
		return !v, true
	}

	if left, right, op, ok := splitComparison(expr); ok {
		lv, lok := resolveExprValue(left, data)
		rv, rok := resolveExprValue(right, data)
		if !lok || !rok {
			return false, false
		}
		eq := fmt.Sprintf("%v", lv) == fmt.Sprintf("%v", rv)
		if op == "==" {
			return eq, true
		}
		return !eq, true
	}

	if strings.EqualFold(expr, "true") {
		return true, true
	}
	if strings.EqualFold(expr, "false") {
		return false, true
	}

	v, ok := resolveExprValue(expr, data)
	if !ok {
		return false, false
	}
	return truthy(v), true
}

func resolveExprValue(token string, data map[string]any) (any, bool) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, false
	}
	if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") && len(token) >= 2 {
		return token[1 : len(token)-1], true
	}
	if strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'") && len(token) >= 2 {
		return token[1 : len(token)-1], true
	}
	if strings.EqualFold(token, "true") {
		return true, true
	}
	if strings.EqualFold(token, "false") {
		return false, true
	}
	if f, err := strconv.ParseFloat(token, 64); err == nil {
		return f, true
	}
	if !isIdentifierPath(token) {
		return token, true
	}

	parts := strings.Split(token, ".")
	var cur any = data
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return token, true
		}
		cur, ok = m[p]
		if !ok {
			return token, true
		}
	}
	return cur, true
}

func truthy(v any) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case string:
		s := strings.TrimSpace(strings.ToLower(t))
		return s != "" && s != "false" && s != "0" && s != "no" && s != "off"
	case int:
		return t != 0
	case int64:
		return t != 0
	case float64:
		return t != 0
	default:
		return true
	}
}

func splitComparison(expr string) (string, string, string, bool) {
	ops := []string{"==", "!="}
	depth := 0
	inQuote := rune(0)
	r := []rune(expr)
	for i := 0; i < len(r)-1; i++ {
		ch := r[i]
		if inQuote != 0 {
			if ch == inQuote {
				inQuote = 0
			}
			continue
		}
		if ch == '\'' || ch == '"' {
			inQuote = ch
			continue
		}
		if ch == '(' {
			depth++
			continue
		}
		if ch == ')' {
			depth--
			continue
		}
		if depth != 0 {
			continue
		}
		for _, op := range ops {
			o := []rune(op)
			if r[i] == o[0] && r[i+1] == o[1] {
				left := strings.TrimSpace(string(r[:i]))
				right := strings.TrimSpace(string(r[i+2:]))
				if left == "" || right == "" {
					return "", "", "", false
				}
				return left, right, op, true
			}
		}
	}
	return "", "", "", false
}

func splitTopLevel(expr, op string) []string {
	depth := 0
	inQuote := rune(0)
	r := []rune(expr)
	opRunes := []rune(op)
	var parts []string
	start := 0

	for i := 0; i <= len(r)-len(opRunes); i++ {
		ch := r[i]
		if inQuote != 0 {
			if ch == inQuote {
				inQuote = 0
			}
			continue
		}
		if ch == '\'' || ch == '"' {
			inQuote = ch
			continue
		}
		if ch == '(' {
			depth++
			continue
		}
		if ch == ')' {
			depth--
			continue
		}
		if depth != 0 {
			continue
		}

		matched := true
		for j := 0; j < len(opRunes); j++ {
			if r[i+j] != opRunes[j] {
				matched = false
				break
			}
		}
		if matched {
			parts = append(parts, strings.TrimSpace(string(r[start:i])))
			start = i + len(opRunes)
			i += len(opRunes) - 1
		}
	}

	if len(parts) == 0 {
		return []string{strings.TrimSpace(expr)}
	}
	parts = append(parts, strings.TrimSpace(string(r[start:])))
	return parts
}

func trimBalancedParens(s string) string {
	for {
		s = strings.TrimSpace(s)
		if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
			return s
		}
		depth := 0
		balanced := true
		for i, ch := range s {
			if ch == '(' {
				depth++
			} else if ch == ')' {
				depth--
				if depth == 0 && i != len(s)-1 {
					balanced = false
					break
				}
			}
			if depth < 0 {
				balanced = false
				break
			}
		}
		if !balanced || depth != 0 {
			return s
		}
		s = s[1 : len(s)-1]
	}
}

func isIdentifierPath(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if !(unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.') {
			return false
		}
	}
	return true
}
