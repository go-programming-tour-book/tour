package sql2strcut

import (
	"os"
	"text/template"

	"github.com/go-programming-tour-book/tour/internal/word"
)

const strcutTpl = `type {{.TableName | UnderscoreToUpperCamelCase}}Model struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $tLen := len .Type }} {{ if gt $tLen 0 }}{{.Name | UnderscoreToUpperCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | UnderscoreToUpperCamelCase}}) TableName() string {
	return "{{.TableName}}"
}`

type StructTemplate struct {
	strcutTpl string
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{strcutTpl: strcutTpl}
}

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType],
			Tag:     "`json:\"" + column.ColumnName + "\"`",
			Comment: column.ColumnComment,
		})
	}
	return tplColumns
}

func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"UnderscoreToUpperCamelCase": word.UnderscoreToUpperCamelCase,
	}).Parse(t.strcutTpl))

	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns:   tplColumns,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}

	return nil
}
