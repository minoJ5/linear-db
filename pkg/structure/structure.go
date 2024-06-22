package structure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Databases struct {
	Databases []*Database   `json:"databases"`
	WriteLock *sync.RWMutex `json:"-"`
}

type Database struct {
	Index     int           `json:"index"`
	Name      string        `json:"name"`
	Tables    []*Table      `json:"tables"`
	WriteLock *sync.RWMutex `json:"-"`
}

type DatabaseIdentifier struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// TODO: See if this is really needed
func (ds *Databases) DatabaseExists(db *Database) bool {
	for _, d := range ds.Databases {
		if d.Name == db.Name {
			return true
		}
	}
	return false
}

// TODO: Maybe change input to a pointer to a database
// Maps do not release memory, so probably this is the most
// efficient way to do this.
func (ds *Databases) IndexOfDatabase(m string) int {
	for i, d := range ds.Databases {
		if m == d.Name {
			return i
		}
	}
	return -1
}

func (ds *Databases) FindDatabase(name string) (*Database, error) {
	ds.WriteLock.RLock()
	defer ds.WriteLock.RUnlock()
	databaseIndex := ds.IndexOfDatabase(name)
	if databaseIndex == -1 {
		return nil, fmt.Errorf("database %s does not exist", name)
	}
	return ds.Databases[databaseIndex], nil
}

func (d *Database) AppendDatabaseResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Database [%s] created.", d.Name)
}

func (ds *Databases) AppendDatabase(d *Database, w http.ResponseWriter) error {
	ds.WriteLock.Lock()
	defer ds.WriteLock.Unlock()
	d.Index = len(ds.Databases)
	if ds.IndexOfDatabase(d.Name) != -1 {
		return fmt.Errorf("database %s already exists", d.Name)
	}
	ds.Databases = append(ds.Databases, d)
	return nil
}

func (ds *Databases) DeleteDatabase(name string) error {
	ds.WriteLock.Lock()
	defer ds.WriteLock.Unlock()
	i := ds.IndexOfDatabase(name)
	if i == -1 {
		return fmt.Errorf("database %s does not exist", name)
	}
	ds.Databases[i] = ds.Databases[len(ds.Databases)-1]
	ds.Databases[len(ds.Databases)-1] = nil
	ds.Databases = ds.Databases[:len(ds.Databases)-1]
	if i < len(ds.Databases){
		ds.Databases[i].Index = i
	}
	return nil
}

func (d *DatabaseIdentifier) DeleteDatabaseResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Database [%s] deleted: `%+v`:", d.Name, *d)
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

// TODO: Use this instead of ReadBodyCreateDatabase
func (d *DatabaseIdentifier) ReadBodyDatabaseGeneral(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(d)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}

type TableQuery struct {
	Database string `json:"database_name"`
	Table    string `json:"table_name"`
}

type DatabaseTable struct {
	DatabaseIndex int    `json:"index"`
	Table         *Table `json:"table"`
}

type DatabasesTables struct {
	Tables    []*DatabaseTable `json:"tables"`
	WriteLock *sync.RWMutex    `json:"-"`
}

type Table struct {
	Index         int           `json:"index"`
	Name          string        `json:"name"`
	DatabaseIndex int           `json:"database_index"`
	Database      string        `json:"database_name"`
	Columns       []*Column     `json:"columns"`
	WriteLock     *sync.RWMutex `json:"-"`
}

type TableIdentifier struct {
	Database string `json:"database_name"`
	Index    int    `json:"index"`
	Name     string `json:"name"`
}

type Column struct {
	Index     int           `json:"index"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Values    []interface{} `json:"values"`
	WriteLock *sync.RWMutex `json:"-"`
}

func GetTalbes(ds *Databases, dn string, tn string) ([]*Table, error) {
	ds.WriteLock.RLock()
	defer ds.WriteLock.RUnlock()
	var dbIndex int = ds.IndexOfDatabase(dn)
	if dbIndex == -1 {
		return nil, fmt.Errorf("databse %s does not exits", dn)
	}
	ts := make([]*Table, 0)
	if tn == "*" {
		ts = append(ts, ds.Databases[dbIndex].Tables...)
	} else {
		var tIndex = ds.TableIndex(dbIndex, tn)
		if tIndex == -1 {
			return nil, fmt.Errorf("table %s does not exits in database %s", tn, dn)
		}
		ts = append(ts, ds.Databases[dbIndex].Tables[tIndex])
	}
	return ts, nil
}

func (ds *Databases) TableExists(i int, t *Table) bool {
	for _, d := range ds.Databases[i].Tables {
		if d.Name == t.Name {
			return true
		}
	}
	return false
}

func (ds *Databases) TableIndex(i int, tn string) int {
	for ti, d := range ds.Databases[i].Tables {
		if d.Name == tn {
			return ti
		}
	}
	return -1
}

func (d *Database) TableIndexDatabase(name string) int {
	for i, t := range d.Tables {
		if name == t.Name {
			return i
		}
	}
	return -1
}

func (t *Table) AppendTableResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Table [%s] created\n%+v", t.Name, *t)
}

func (ds *Databases) AppendTable(t *Table) error {
	// TODO: Don't lock all databases, just the one where the table is inserted
	ds.WriteLock.Lock()
	defer ds.WriteLock.Unlock()
	var dbIndex int = ds.IndexOfDatabase(t.Database)
	var db *Database = ds.Databases[dbIndex]
	t.DatabaseIndex = dbIndex
	t.Index = len(db.Tables)
	if dbIndex == -1 {
		return fmt.Errorf("databse %s does not exits", db.Name)
	}
	if ds.TableExists(dbIndex, t) {
		return fmt.Errorf("table %s alreay exits in database %s", t.Name, db.Name)
	}
	db.Tables = append(db.Tables, t)
	t = nil
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

func (tq *TableQuery) ReadBodyGetTables(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(tq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}

// TODO: Use this instead of ReadBodyCreateTable
func (t *TableIdentifier) ReadBodyTableGeneral(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}

func (d *Database) DeleteTable(name string) error {
	// d.WriteLock.Lock()
	// defer d.WriteLock.Unlock()
	tableIndex := d.TableIndexDatabase(name)
	if tableIndex == -1 {
		return fmt.Errorf("table %s does not exits in database %s", name, d.Name)
	}
	d.Tables[tableIndex] = d.Tables[len(d.Tables)-1]
	d.Tables[len(d.Tables)-1] = nil
	d.Tables = d.Tables[:len(d.Tables)-1]
	if tableIndex < len(d.Tables) {
		d.Tables[tableIndex].Index = tableIndex
	}
	return nil
}

func (t *TableIdentifier) DeleteTableResponse(d *Database, w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Table [%s] in Database [%s] deleted.", t.Name, d.Name)
}
