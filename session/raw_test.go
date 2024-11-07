package session

import (
	"database/sql"
	dialect "geeORM/Dalect"
	"os"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var TestDB *sql.DB
var sqlite3Dialect dialect.Dialect

// test.会在测试单元执行前先被执行
func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../gee.db")
	sqlite3Dialect, _ = dialect.GetDialect("sqlite3")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, sqlite3Dialect)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRow(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	row := s.Raw("SELECT count(*) from User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 2 {
		t.Fatal("failed to query db")
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	row, err := s.Raw("SELECT Name from User").Query()
	if err != nil {
		t.Fatal("failed to query")
	}
	var names []string = make([]string, 0)
	var name string
	expectNames := []string{"Tom", "Sam"}
	for row.Next() {
		if err = row.Scan(&name); err != nil {
			t.Fatal("failed to Scan ")
		}
		names = append(names, name)
	}
	if !reflect.DeepEqual(names, expectNames) {
		t.Fatal("failed to Scan all names")
	}
}
