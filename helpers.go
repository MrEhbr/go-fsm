package fsm

import (
	"strings"
	"unicode"
)

// ToCamelCase transform string to upper or lower camel case
func ToCamelCase(s string, upper bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var (
		nextUpper = upper
		prev      rune
		builder   strings.Builder
	)

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			nextUpper = true
			continue
		}

		if nextUpper {
			builder.WriteRune(unicode.ToUpper(r))
			nextUpper = false
			continue
		}

		if unicode.IsLower(prev) {
			builder.WriteRune(r)
			continue
		}

		builder.WriteRune(unicode.ToLower(r))
		prev = r
	}

	return builder.String()
}

func ToSnackCase(s string, upper bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	adjustCase := unicode.ToLower
	if upper {
		adjustCase = unicode.ToUpper
	}

	var (
		builder strings.Builder
		prev    rune
		curr    rune
	)
	for _, next := range s {
		switch {
		case !unicode.IsLetter(curr) && !unicode.IsDigit(curr):
			if unicode.IsLetter(prev) || unicode.IsDigit(prev) {
				builder.WriteRune('_')
			}
		case unicode.IsUpper(curr):
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && unicode.IsLower(next)) {
				builder.WriteRune('_')
			}
			builder.WriteRune(adjustCase(curr))
		default:
			builder.WriteRune(adjustCase(curr))
		}

		prev = curr
		curr = next
	}

	if unicode.IsUpper(curr) && unicode.IsLower(prev) && prev != 0 {
		builder.WriteRune('_')
	}
	builder.WriteRune(adjustCase(curr))

	return builder.String()
}

func InStrings(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}

	return false
}
