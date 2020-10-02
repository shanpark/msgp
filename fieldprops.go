package msgp

import (
	"reflect"
	"strings"
)

// FieldProps represents field properties of struct for struct packing.
type FieldProps struct {
	Name      string
	Skip      bool
	OmitEmpty bool
	String    bool
}

func (fp *FieldProps) parseTag(field reflect.StructField) {
	tag := field.Tag.Get("msgp")
	if tag == "" {
		fp.Name = field.Name
	} else {
		name, opts := parseTag(tag)
		if name == "-" {
			if strings.TrimSpace(string(opts)) == "" {
				fp.Skip = true
			} else {
				fp.Name = "_"
				return
			}
		} else if len(name) > 0 {
			fp.Name = name
		} else {
			fp.Name = field.Name
		}

		if opts.Contains("omitempty") {
			fp.OmitEmpty = true
		}

		if opts.Contains("string") {
			fp.String = true
		}
	}
}
