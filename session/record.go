package session

import (
	"errors"
	"geeORM/clause"
	"geeORM/log"
	"reflect"
)

// Insert can insert all values to table which is link to the session model
// * Usage like s.Insert(u1, u2, ...)
// * the elem of the values need be a pointer not the instance of struct
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	table := s.Model(values[0]).RefTable()

	s.clause.Set(clause.INSERT, table.Name)
	for _, val := range values {
		s.CallHook(BeforeInsert, reflect.ValueOf(val).Interface())
		recordValues = append(recordValues, table.RecordValues(val)) // save all values
		s.CallHook(AfterInsert, reflect.ValueOf(val).Interface())
	}

	s.clause.Set(clause.VALUES, recordValues...)
	// generate all values for sql
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	log.Info(sql, vars)
	result, err := s.db.Exec(sql, vars...)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return result.RowsAffected()
}

// Find can save all result int values points
// Find all columns the the conditions set before
// values must be a slice pointer type
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	s.Model(reflect.New(destType).Elem().Interface())
	table := s.RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).Query()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		vals := make([]interface{}, 0, len(table.FieldNames))
		for _, name := range table.FieldNames {
			vals = append(vals, dest.FieldByName(name).Addr().Interface())
		}
		s.CallHook(BeforeQuery, dest.Addr().Interface())
		if err = rows.Scan(vals...); err != nil {
			return err
		}
		s.CallHook(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// support map[string]interface{}
// also kv list
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallHook(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallHook(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.CallHook(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	s.CallHook(AfterDelete, nil)
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var temp int64
	if err := row.Scan(&temp); err != nil {
		return 0, err
	}
	return temp, nil
}

// Limit adds limit condition to clause
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where add limit condition to clause
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// get the first
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
