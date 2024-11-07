package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    // 将struct 的类型转换为数据库的类型
	TableExistSQL(tableName string) (string, []interface{}) // 返回某个表是否存在的SQL语句，参数是表名
}

// each database-dialect will call RegisterDialect in init function
// then we can use the dialect at any time
func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
