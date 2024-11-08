package session

import (
	"geeORM/log"
	"reflect"
	"testing"
)

var (
	u1 User = User{"Tom", 12}
	u2 User = User{"LY", 23}
	u3 User = User{"MJT", 2}
)

func TestRecordInsert(t *testing.T) {
	s := NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	num, err := s.Insert(&u1, &u2, &u3)
	if err != nil || num != 3 {
		t.Fatalf("err %v num: %d", err, num)
	}
}

func TestRecordFind(t *testing.T) {
	s := NewSession().Model(&User{})
	var users []User
	s.Find(&users)
	log.Infof("%+v ", users)
	if !reflect.DeepEqual(users, []User{u1, u2, u3}) {
		t.Fatal("not found the aim data")
	}
}
