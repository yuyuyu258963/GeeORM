package session

import (
	"geeORM/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallHook calls the registered hooks
func (s *Session) CallHook(method string, arg interface{}) {
	f := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if arg != nil {
		f = reflect.ValueOf(arg).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if f.IsValid() {
		if v := f.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
