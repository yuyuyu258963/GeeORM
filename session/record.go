package session

import (
	"geeORM/clause"
	"geeORM/log"
	"reflect"
)

// Insert can insert all values to table which is link to the session model
// * Usage like s.Insert(u1, u2, ...)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	table := s.Model(values[0]).RefTable()
	s.clause.Set(clause.INSERT, table.Name)
	for _, val := range values {
		recordValues = append(recordValues, table.RecordValues(val)) // save all values
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
		if err = rows.Scan(vals...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
