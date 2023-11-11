package dbscan

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

const (
	debug = false
)

// Into will create a new instance of T, scan the row into it, and return it. T must be some kind of struct.
func Into[T any](row *sql.Rows) (T, error) {
	var into T
	if debug {
		intoType := reflect.TypeOf(into)
		if intoType.Kind() != reflect.Struct {
			log.Fatalf("dbscan.Into can only scan into structs, but has a type parameter of type %T", into)
		}
	}
	intoValue := reflect.ValueOf(&into).Elem()
	intoPtrs := make([]interface{}, intoValue.NumField())
	for i := 0; i < intoValue.NumField(); i++ {
		intoPtrs[i] = intoValue.Field(i).Addr().Interface()
	}

	err := row.Scan(intoPtrs...)
	if err != nil {
		return into, fmt.Errorf("dbscan.Into failed to scan into a %T: %v", into, err)
	}

	return into, nil
}

// All will scan all the rows into a slice of the given type T (must be a struct).
func All[T any](rows *sql.Rows) ([]T, error) {
	var all []T
	for rows.Next() {
		t, err := Into[T](rows)
		if err != nil {
			return all, fmt.Errorf("failed to process query result rows: %v", err)
		}
		all = append(all, t)
	}

	return all, nil
}
