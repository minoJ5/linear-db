package httpserver

import (
	"encoding/json"
	"fmt"
	"linear-db/pkg/misc"
	sr "linear-db/pkg/structure"
	"net/http"
	"regexp"
	"time"
)

var (
	ds *sr.Databases
	dt *sr.DatabasesTables
)

func init() {
	ds = &sr.Databases{
		Databases: make([]sr.Database, 0),
	}
	dt = &sr.DatabasesTables{
		Tables: make([]sr.DatabaseTables, 0),
	}
}

func live(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server is live: Now %v\n", time.Now())))
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method [%s] is not allowed!\n", r.Method), http.StatusMethodNotAllowed)
}

func readBodyCreateDatabase(d *sr.Database, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(d)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}

func createLDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		methodNotAllowedResponse(w, r)
		return
	}
	d := new(sr.Database)
	err := readBodyCreateDatabase(d, w, r)
	if err != nil {
		//fmt.Println(fmt.Errorf("read body err:\n%s", err))
		return
	}
	if len(d.Name) == 0 {
		http.Error(w, "Database's name cannot be empty", http.StatusConflict)
		return
	}
	if match, _ := regexp.MatchString("^[a-zA-Z0-9-_]+$", d.Name); !match {
		http.Error(w, fmt.Sprintf("Database's name [%s] not valid, allowed are names with characters from A(a) to Z(z) with - or _ as seperation!", d.Name), http.StatusConflict)
		return
	}
	if misc.DatabaseExists(ds, d.Name) {
		http.Error(w, fmt.Sprintf("Database [%s] already exitsts", d.Name), http.StatusConflict)
		return
	}
	d.Index = len(ds.Databases)
	ds.WriteLock.RLock()
	defer ds.WriteLock.RUnlock()
	ds.Databases = append(ds.Databases, *d)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Database [%s] created: `%+v`:", d.Name, *d)
}

func listLDatabases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET")
	if r.Method != http.MethodGet {
		methodNotAllowedResponse(w, r)
		return
	}
	resp, err := json.Marshal(ds.Databases)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding error:\n%+v\n", ds), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", string(resp))
}

func cselect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET")
	if r.Method != http.MethodGet {
		methodNotAllowedResponse(w, r)
		return
	}
	w.Write([]byte("Get ok!"))
}

func readBodyCreateTable(t *sr.Table, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoing JSON:\n%s\n", err), http.StatusBadRequest)
		return err
	}
	return nil
}
func createTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowedResponse(w, r)
		return
	}
	t := new(sr.Table)

	err := readBodyCreateTable(t, w, r)
	if err != nil {
		return
	}

	tt := new(sr.DatabaseTables)
	ix := misc.IndexOf(t.Database, ds)
	if ix == -1 {
		http.Error(w, fmt.Sprintf("Error: Database %s does not exist", t.Database), http.StatusConflict)
		return
	}
	tt.LDatabaseIndex = ix
	tt.Tables = append(tt.Tables, *t)
	dt.Tables = append(dt.Tables, *tt)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Table [%s] created\n%+v", t.Name, *t)
	t = nil
	tt = nil
}

func listTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET")
	if r.Method != http.MethodGet {
		methodNotAllowedResponse(w, r)
		return
	}
	resp, err := json.Marshal(dt)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding error:\n%+v\n", ds), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", string(resp))
}
