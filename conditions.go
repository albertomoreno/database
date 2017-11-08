package database

import (
	"fmt"
	"strings"
	"reflect"
)

type Condition interface {
	SQL() string
	Values() []interface{}
}

type simpleCondition struct {
	sql   string
	value interface{}
}

func (cond *simpleCondition) SQL() string {
	if !strings.Contains(cond.sql, " ") {
		return fmt.Sprintf("%s = ?", cond.sql)
	}

	if strings.Contains(cond.sql, " IN") {
		v := reflect.ValueOf(cond.value)
		placeholders := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			placeholders[i] = "?"
		}
		return fmt.Sprintf("%s (%s)", cond.sql, strings.Join(placeholders, ", "))
	}

	if !strings.Contains(cond.sql, "?") {
		return fmt.Sprintf("%s ?", cond.sql)
	}

	return cond.sql
}

func (cond *simpleCondition) Values() []interface{} {
	if strings.Contains(cond.sql, " IN") {
		v := reflect.ValueOf(cond.value)
		var values []interface{}
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i).Interface())
		}
		return values
	}

	return []interface{}{cond.value}
}
