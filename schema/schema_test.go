package schema

import (
	dialect "geeORM/Dalect"
	"testing"
)

type User struct {
	Id   int `geeorm:"primary key"`
	Name string
	age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{Id: 1}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Id").Tag != "primary key" {
		t.Fatal("failed to parse tag")
	}
}
