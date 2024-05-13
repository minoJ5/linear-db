package httpserver

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
	//cmw "linear-db/pkg/middleware"
)

func HttpServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Use(cmw.Memory)
	url := "127.0.0.1:8001"
	log.Printf("server started on %s\n", url)
	r.Get("/live", live)
	r.HandleFunc("/createdb", createLDatabase)
	r.HandleFunc("/select", cselect)
	r.HandleFunc("/listdbs", listLDatabases)
	r.HandleFunc("/createtable", createTable)
	r.HandleFunc("/listtables", listTables)
	err := http.ListenAndServe(url, r)
	log.Fatalf("server crashed:\n %s\n", err)
}
