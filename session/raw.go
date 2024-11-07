package session

import (
	"database/sql"
	"strings"

	dialect "geeORM/Dalect"
	"geeORM/log"
	"geeORM/schema"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	// 用于拼接生成SQL语句和SQL语句中占位符的对应值
	sql     strings.Builder
	sqlVars []interface{}
}

// *sql.DB是通过sql.Open()方法连接数据库成功后返回的指针
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// 清空会话
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// log info level message
func (s *Session) logInfo() {
	log.Info(s.sql.String(), s.sqlVars)
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	s.logInfo()
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(s.sql.String(), s.sqlVars)
	}
	return
}

func (s *Session) Query() (row *sql.Rows, err error) {
	defer s.Clear()
	s.logInfo()
	if row, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(s.sql.String(), s.sqlVars)
	}
	return
}

// only get one row if exits
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	s.logInfo()
	row := s.DB().QueryRow(s.sql.String(), s.sqlVars...)
	return row
}
