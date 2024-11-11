package geeorm

import (
	"database/sql"
	"fmt"
	dialect "geeORM/Dalect"
	"geeORM/log"
	"geeORM/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

type TxFunc func(*session.Session) (interface{}, error)

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

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback()
		} else { // 如果过程中没有error发生那就可以Commit
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()

			err = s.Commit() // err is nil; if commit returns error update err
		}
	}()

	return f(s)
}

// diff = a - b
func difference(a []string, b []string) (diff []string) {
	mpB := make(map[string]bool)
	for _, i := range b {
		mpB[i] = true
	}
	for _, v := range a {
		if _, ok := mpB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate Table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s not exists", s.RefTable().Name)
			return
		}
		table := s.RefTable()
		row, _ := s.Raw(fmt.Sprintf("SELECT * from %s limit 1;", table.Name)).Query()
		columns, _ := row.Columns()

		delColumns := difference(columns, table.FieldNames)
		addColumns := difference(table.FieldNames, columns)
		log.Infof("delete columns %v, add columns %v", delColumns, addColumns)
		// 先处理新增的列 使用 ALTER
		for _, column_name := range addColumns {
			f := s.RefTable().GetField(column_name)
			sqlStr := fmt.Sprintf("Alter Table %s add column %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delColumns) == 0 {
			return
		}
		// 处理删除掉的列
		tmp := "tmp_" + table.Name
		fields := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fields, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s ;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
