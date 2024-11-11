package session

import (
	"database/sql"
	"strings"

	dialect "geeORM/Dalect"
	"geeORM/clause"
	"geeORM/log"
	"geeORM/schema"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	// 用于拼接生成SQL语句和SQL语句中占位符的对应值
	sql     strings.Builder
	sqlVars []interface{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

var _ CommonDB = &sql.Tx{} // to insure *sql.Tx has implement the CommonDB interface

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
	s.clause = clause.Clause{}
	s.sqlVars = nil
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// 设置原始的sql语句以及其中需要填充的值
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
