package session

import (
	"geeORM/log"
	"reflect"
	"testing"
)

var (
	u1 User = User{"Tom", 2}
	u2 User = User{"LY", 23}
	u3 User = User{"MJT", 2}
	u4 User = User{"Ywh", 2}
)

func TestRecordInsert(t *testing.T) {
	s := NewSession().Model(&User{})
	s.DropTable()
	s.CreateTable()
	num, err := s.Insert(&u1, &u2, &u3, &u4)
	if err != nil || num != 4 {
		t.Fatalf("err %v num: %d", err, num)
	}
}

func TestRecordFind(t *testing.T) {
	s := NewSession().Model(&User{})
	var users []User
	s.Find(&users)
	log.Infof("%+v ", users)
	if !reflect.DeepEqual(users, []User{u1, u2, u3, u4}) {
		t.Fatal("not found the aim data")
	}
}

func TestRecordUpdate(t *testing.T) {
	s := NewSession().Model(&User{})
	m := map[string]interface{}{
		"Age": 2,
	}
	if n, err := s.Where("Name = ?", "Tom").Update(m); err != nil || n != 1 {
		t.Fatal(err, n)
	}
}

func TestRecordFirst(t *testing.T) {
	s := NewSession().Model(&User{})
	user := &User{}
	s.Where("Age = ?", 2).First(user)
	log.Infof("%+v\n", user)
	b := &User{}
	if s.Where("Age = ?", 9999).First(b) == nil {
		t.Fatal("error found one")
	}
}

// 测试Limit 和 orderby
func TestRecordOrderBy(t *testing.T) {
	s := NewSession().Model(&User{})
	var users []User
	s.OrderBy("Age desc").Limit(2).Find(&users)
	log.Infof("%+v\n", users)
	if len(users) != 2 {
		t.Fatal("error found two")
	}
}

func TestDelete(t *testing.T) {
	s := NewSession().Model(&User{})
	n, err := s.Where("age > ?", 2).Limit(1).Delete()
	if err != nil || n != 1 {
		t.Fatal(n, err)
	}
}

// Count功能是否正常
func TestCount(t *testing.T) {
	s := NewSession().Model(&User{})
	n, err := s.Where("age < ?", 10).Count()
	if n != 3 || err != nil {
		t.Fatal(n, err)
	}
}
