package schema

import (
	dialect "geeORM/Dalect"
	"reflect"
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

func TestRecordValues(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	vals := schema.RecordValues(&User{Id: 1, Name: "Join"})
	if !reflect.DeepEqual(vals, []interface{}{1, "Join"}) {
		t.Fatalf("failed to get record values %+v", vals)
	}
}
