package msgp

import (
	"strings"
)

// tagOptions is the string following a comma in a struct field's "msgp"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's msgp tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return strings.TrimSpace(tag[:idx]), tagOptions(tag[idx+1:])
	}
	return strings.TrimSpace(tag), tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := strings.TrimSpace(string(o))
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
