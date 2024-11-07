package schema

import (
	dialect "geeORM/Dalect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string // 所有字段名、即列名
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 传入一个结构体指针，通过Dialect映射为对应数据库的表结构体
func Parse(v interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(v)).Type()
	schema := &Schema{
		Model:    modelType,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := range modelType.NumField() {
		item := modelType.Field(i)
		// 只考虑struct中的非嵌入成员变量；以及对外暴露的成员变量
		if !item.Anonymous && ast.IsExported(item.Name) {
			field := &Field{
				Name: item.Name,
				// reflect.New返回一个指向该Type的指针， reflect.Indirect返回指针指向的值
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(item.Type))),
			}
			if tag, ok := item.Tag.Lookup("geeorm"); ok {
				field.Tag = tag
			}
			schema.fieldMap[field.Name] = field
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
		}
	}
	return schema
}
