package geeorm

import (
	"geeORM/session"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	// Helper 标记函数作为测试辅助函数，使得在这段代码中发生的仍和失败都被
	// 报告为发生在这段代码的开始处，而不是实际发生的位置
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", transactionRollback)
	t.Run("commit", transactionCommit)
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(
		func(s *session.Session) (result interface{}, err error) {
			_ = s.Model(&User{}).CreateTable()
			_, err = s.Insert(&User{"Tom", 18})
			_, err = s.Insert(&User{"Jack", 8})
			return
		},
	)

	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()
	_, err := engine.Transaction(
		func(s *session.Session) (result interface{}, err error) {
			_, err = s.Insert(&User{"Tom", 18})
			return
		},
	)

	u := &User{}
	_ = s.First(u)
	if err != nil || !s.HasTable() || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func TestEngineMigrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, a integer);").Exec()
	_, _ = s.Raw("INSERT INTO User('Name') values (?), (?), (?)", "Tom", "JJ", "Jit").Exec()
	engine.Migrate(&User{}) // 测试User表迁移

	rows, _ := s.Raw("SELECT * from User").Query()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
