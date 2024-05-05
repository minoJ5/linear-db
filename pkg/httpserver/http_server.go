package httpserver

import (
	"log"
	"net/http"
)

func HttpServer() {
	mux := http.NewServeMux()
	url := "127.0.0.1:8001"
	log.Printf("server started on %s\n", url)
	mux.HandleFunc("/live", live)
	mux.HandleFunc("/createdb", createLDatabase)
	mux.HandleFunc("/select", cselect)
	mux.HandleFunc("/listdbs", listLDatabases)
	mux.HandleFunc("/createtable", createTable)
	mux.HandleFunc("/listtables", listTables)
	err := http.ListenAndServe(url, mux)
	log.Fatalf("server crashed:\n %s\n", err)
}
