package structure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Databases struct {
	Databases []Database
	WriteLock sync.RWMutex
}

type Database struct {
	Index int
	Name  string
}

func (ds *Databases) DatabaseExists(db *Database) bool {
	for _, d := range ds.Databases {
		if d.Name == db.Name {
			return true
		}
	}
	return false
}

func (ds *Databases) IndexOf(m string) int {
	for i, d := range ds.Databases {
		if m == d.Name {
			return i
		}
	}
	return -1
}

func (d *Database) AppendDatabaseResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Database [%s] created: `%+v`:", d.Name, *d)
}

func (ds *Databases) AppendDatabase(d *Database, w http.ResponseWriter) {
	d.Index = len(ds.Databases)
	ds.WriteLock.RLock()
	defer ds.WriteLock.RUnlock()
	ds.Databases = append(ds.Databases, *d)
}

func (d *Database) ReadBodyCreateDatabase(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(d)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}

type DatabaseTable struct {
	DatabaseIndex int
	Table         Table
}

type DatabasesTables struct {
	Tables    []DatabaseTable
	WriteLock sync.RWMutex
}

type Table struct {
	Index         int
	Name          string
	DatabaseIndex int
	Database      string
	Columns       map[string]Column
}

func (dt *DatabasesTables) TableExists(t *Table) bool {
	for _, d := range dt.Tables {
		if d.Table.Name == t.Name && d.DatabaseIndex == t.DatabaseIndex {
			return true
		}
	}
	return false
}

func (t *Table) AppendTableResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Table [%s] created\n%+v", t.Name, *t)
}

func (dt *DatabasesTables) AppendTable(t *Table, ds *Databases) error{
	ds.WriteLock.RLock()
	dt.WriteLock.RLock()
	defer dt.WriteLock.RUnlock()
	defer ds.WriteLock.RUnlock()
	td := new(DatabaseTable)
	t.Index = 0
	t.DatabaseIndex = ds.IndexOf(t.Database)
	td.DatabaseIndex = t.DatabaseIndex
	if dt.TableExists(t) {
		return fmt.Errorf("table %s alreay exits in database %s", t.Name, t.Database)
	}
	td.Table = *t
	dt.Tables = append(dt.Tables, *td)
	t = nil
	td = nil
	return nil
}

func (t *Table) ReadBodyCreateTable(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}

type Column struct {
	Index  int
	Name   string
	Type   string
	Values []interface{}
}
