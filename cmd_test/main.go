package main

import (
	"fmt"
	geeorm "geeORM"
	"geeORM/log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func serveSqlite() {
	e, err := geeorm.NewEngine("sqlite3", "gee.db") // database driver name , dsn
	if err != nil {
		log.Error(err)
	}
	defer func() {
		e.Close()
	}()

	s := e.NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec() // 测试错误日志是否能正常输出

	result, err := s.Raw("INSERT INTO User('name') values (?),(?)",
		"Tom", "Sam").Exec()
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Info(affected)
	}

	row := s.Raw("SELECT Name FROM User LIMIT 2").QueryRow()
	var name string
	if err := row.Scan(&name); err == nil {
		// log.Info(name)
		fmt.Println(name)
	}
}

type User struct {
	Id   int `geeorm:"PRIMARY KEY" json:"id"`
	Name string
	Age  int
}

func main() {
	var u *User = &User{Id: 1, Name: "Tom", Age: 1}
	t := reflect.Indirect(reflect.ValueOf(u)).Type()
	fmt.Println(reflect.TypeOf(u))
	for i := range t.NumField() {
		fieldItem := t.Field(i)
		fmt.Println(fieldItem.Name)
		fmt.Println(fieldItem.Type)
		fmt.Println(fieldItem.Tag.Get("geeorm"))
	}
}
