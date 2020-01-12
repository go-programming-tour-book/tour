package json2struct

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-programming-tour-book/tour/internal/word"
)

const (
	TYPE_MAP_STRING_INTERFACE = "map[string]interface {}"
	TYPE_INTERFACE            = "[]interface {}"
)

type Parser struct {
	Source     map[string]interface{}
	Output     Output
	Children   Output
	StructTag  string
	StructName string
}

type Output []string

func (o *Output) appendSegment(format, title string, args ...interface{}) {
	s := []interface{}{}
	s = append(s, word.UnderscoreToUpperCamelCase(title))
	if len(args) != 0 {
		s = append(s, args...)
		format = "    " + format
	}
	*o = append(*o, fmt.Sprintf(format, s...))
}

func (o *Output) appendSuffix() {
	*o = append(*o, "}\n")
}

func NewParser(s string) (*Parser, error) {
	source := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &source); err != nil {
		return nil, err
	}
	return &Parser{
		Source:     source,
		StructTag:  "type %s struct {",
		StructName: "tour",
	}, nil
}

func (p *Parser) Json2Struct() string {
	p.Output.appendSegment(p.StructTag, p.StructName)
	for sourceName, sourceValue := range p.Source {
		valueType := reflect.TypeOf(sourceValue).String()
		if valueType == TYPE_INTERFACE {
			p.toParentList(sourceName, sourceValue.([]interface{}), true)
		} else {
			var fields Fields
			fields.appendSegment(sourceName, FieldSegment{
				Format: "%s",
				FieldValues: []FieldValue{
					{CamelCase: false, Value: valueType},
				},
			})
			p.Output.appendSegment("%s %s", fields[0].Name, fields[0].Type)
		}
	}
	p.Output.appendSuffix()
	return strings.Join(append(p.Output, p.Children...), "\n")
}

func (p *Parser) toChildrenStruct(parentName string, values interface{}) {
	p.Children.appendSegment(p.StructTag, parentName)
	for fieldName, fieldValue := range values.(map[string]interface{}) {
		p.Children.appendSegment("%s %s", fieldName, reflect.TypeOf(fieldValue).String())
	}
	p.Children.appendSuffix()
}

func (p *Parser) toParentList(parentName string, values []interface{}, isParent bool) {
	var fields Fields
	for _, v := range values {
		valueType := reflect.TypeOf(v).String()
		if valueType == TYPE_MAP_STRING_INTERFACE {
			for fieldName, fieldValues := range v.(map[string]interface{}) {
				var (
					fieldValueType = reflect.TypeOf(fieldValues).String()
					fieldSegment   = FieldSegment{
						Format: "%s",
						FieldValues: []FieldValue{
							{CamelCase: true, Value: fieldValueType},
						},
					}
				)
				switch fieldValueType {
				case TYPE_INTERFACE:
					p.toParentList(fieldName, fieldValues.([]interface{}), false)
					fieldSegment.Format = "%s%s"
					fieldSegment.FieldValues = []FieldValue{
						{CamelCase: false, Value: "[]"},
						{CamelCase: true, Value: fieldName},
					}
				case TYPE_MAP_STRING_INTERFACE:
					p.toChildrenStruct(fieldName, fieldValues)
					fieldSegment.Format = "%s"
					fieldSegment.FieldValues = []FieldValue{
						{CamelCase: true, Value: fieldName},
					}
				}

				fields.appendSegment(fieldName, fieldSegment)
			}

			p.Children.appendSegment(p.StructTag, parentName)
			for _, field := range fields.removeDuplicate() {
				p.Children.appendSegment("%s %s", field.Name, field.Type)
			}
			p.Children.appendSuffix()
			if isParent {
				valueType = word.UnderscoreToUpperCamelCase(parentName)
			}
		}

		if isParent {
			p.Output.appendSegment("%s %s%s", parentName, "[]", valueType)
		}
		break
	}
}
