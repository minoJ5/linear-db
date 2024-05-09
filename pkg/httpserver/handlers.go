package httpserver

import (
	"encoding/json"
	"fmt"
	sr "linear-db/pkg/structure"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var (
	ds *sr.Databases
	dt *sr.DatabasesTables
)

func init() {
	ds = &sr.Databases{
		Databases: make([]sr.Database, 0),
		WriteLock: new(sync.RWMutex),
	}
	dt = &sr.DatabasesTables{
		Tables: make([]sr.DatabaseTable, 0),
	}
}

func cselect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET")
	if r.Method != http.MethodGet {
		methodNotAllowedResponse(w, r)
		return
	}
	w.Write([]byte("Get ok!"))
}

func live(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server is live: Now %v\n", time.Now())))
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Method [%s] is not allowed!\n", r.Method), http.StatusMethodNotAllowed)
}

func createLDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		methodNotAllowedResponse(w, r)
		return
	}
	d := new(sr.Database)

	err := d.ReadBodyCreateDatabase(w, r)
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
	if ds.DatabaseExists(d) {
		http.Error(w, fmt.Sprintf("Database [%s] already exitsts", d.Name), http.StatusConflict)
		return
	}
	ds.AppendDatabase(d, w)
	d.AppendDatabaseResponse(w)
}

func listLDatabases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET")
	if r.Method != http.MethodGet {
		methodNotAllowedResponse(w, r)
		return
	}
	ds.WriteLock.RLock()
	defer ds.WriteLock.RUnlock()
	resp, err := json.Marshal(ds.Databases)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding error:\n%+v\n", ds), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", string(resp))
}

func createTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowedResponse(w, r)
		return
	}
	t := new(sr.Table)
	err := t.ReadBodyCreateTable(w, r)
	if err != nil {
		return
	}
	if ix := ds.IndexOf(t.Database); ix == -1 {
		http.Error(w, fmt.Sprintf("Error: Database [%s] does not exist", t.Database), http.StatusConflict)
		return
	}
	err = ds.AppendTable(t)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: Table [%s] aready exists in [%s] ", t.Name, t.Database), http.StatusConflict)
		fmt.Println(err)
		return
	}
	t.AppendTableResponse(w)
}

func listTables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST")
	if r.Method != http.MethodPost {
		methodNotAllowedResponse(w, r)
		return
	}
	tq := new(sr.TableQuery)
	err := tq.ReadBodyGetTables(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query invalid:\n%+v\n", tq), http.StatusInternalServerError)
		return
	}
	tables, err := sr.GetTalbes(ds, tq.Database, tq.Table)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query invalid 2:\n%+v\n%s\n", tq, err), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(tables)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding error:\n%+v\n", dt), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", string(resp))
}
