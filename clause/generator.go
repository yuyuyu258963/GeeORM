package clause

import (
	"bytes"
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// 生成指定长度的 ?,?,....
func genBindVars(num int) string {
	var builder bytes.Buffer
	builder.Grow(num*2 - 1)
	for i := 0; i < num; i++ {
		builder.WriteString("?")
		if i != num-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}

// INSERT INTO tableName values (...,...)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	return fmt.Sprintf("INSERT INTO %s", tableName), []interface{}{}
}

// LIMIT $num
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", []interface{}{values[0]}
}

// 用于批量插入
// VALUES ($v1), ($v2) ...
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, val := range values {
		v := val.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// SELECT $fields FROM tableNAME
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// WHERE $desc
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %v", desc), vars
}

// ORDER BY $field
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %v", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
