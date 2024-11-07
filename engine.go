package geeorm

import (
	"database/sql"
	dialect "geeORM/Dalect"
	"geeORM/log"
	"geeORM/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driverName, dsn string) (e *Engine, err error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Error(err)
		return
	}

	// Send a Ping to make sure the database is connected
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dialect, ok := dialect.GetDialect(driverName)
	if !ok {
		log.Errorf("not found the dialect: %s", driverName)
		return
	}

	e = &Engine{db, dialect}
	log.Info("success connect to database")
	return
}

// close the database connection
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("success disconnect from database")
}

// 通过当前的连接创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
