package structure

import "sync"

type Databases struct {
	Databases []Database
	WriteLock sync.RWMutex
}

type Database struct {
	Index int
	Name  string
}

type DatabaseTables struct {
	LDatabaseIndex int
	Tables         []Table
}

type DatabasesTables struct {
	Tables []DatabaseTables
}

type Table struct {
	Index         int
	Name          string
	DatabaseIndex int
	Database      string
	Columns       map[string]Column
}

type Column struct {
	Index  int
	Name   string
	Type   string
	Values []interface{}
}
