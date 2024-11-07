package session

import (
	"errors"
	"fmt"
	"geeORM/log"
	"geeORM/schema"
	"reflect"
	"strings"
)

// Model Set the session link to a table
func (s *Session) Model(value interface{}) *Session {
	// 如果这个session上一次操作的table和上一次是相同的则不用再解析一次了
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// 获取本Session关联的Table
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// Session create table with linked struct
func (s *Session) CreateTable() (err error) {
	table := s.RefTable()
	if table == nil {
		return errors.New("not set linked table")
	}
	var columns []string
	for _, filed := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", filed.Name, filed.Type, filed.Tag))
	}
	destFields := strings.Join(columns, ",")
	_, err = s.Raw(fmt.Sprintf("CREATE TABLE %s  (%s);", table.Name, destFields)).Exec()
	return err
}

// drop the table with link to the session model
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// check the table is created or not
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var temp string
	_ = row.Scan(&temp)
	return temp == s.RefTable().Name
}
