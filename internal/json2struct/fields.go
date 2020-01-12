package json2struct

import (
	"fmt"

	"github.com/go-programming-tour-book/tour/internal/word"
)

type FieldSegment struct {
	Format      string
	FieldValues []FieldValue
}

type FieldValue struct {
	CamelCase bool
	Value     string
}

type Field struct {
	Name string
	Type string
}

type Fields []*Field

func (f *Fields) appendSegment(name string, segment FieldSegment) {
	var s []interface{}
	for _, v := range segment.FieldValues {
		value := v.Value
		if v.CamelCase {
			value = word.UnderscoreToUpperCamelCase(v.Value)
		}

		s = append(s, value)
	}
	*f = append(*f, &Field{Name: word.UnderscoreToUpperCamelCase(name), Type: fmt.Sprintf(segment.Format, s...)})
}

func (f *Fields) removeDuplicate() Fields {
	m := make(map[string]bool)
	fields := Fields{}
	for _, entry := range *f {
		if _, value := m[entry.Name]; !value {
			m[entry.Name] = true
			fields = append(fields, entry)
		}
	}
	return fields
}
