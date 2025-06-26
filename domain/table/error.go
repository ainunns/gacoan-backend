package table

import "errors"

var (
	ErrorGetAllTables  = errors.New("failed to get all tables")
	ErrorGetTableByID  = errors.New("failed to get table by id")
	ErrorTableNotFound = errors.New("table not found")
)
