package httpserver

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	//cmw "linear-db/pkg/middleware"
)

func HttpServer() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("HTTP server: Error loading .env file", err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Use(cmw.Memory)
	log.Printf("server started on %s\n", os.Getenv("URL_HTTP"))
	r.Get("/live", live)
	r.HandleFunc("/createdb", createDatabase)
	r.HandleFunc("/deletedb", deleteDatabase)
	r.HandleFunc("/select", cselect)
	r.HandleFunc("/listdbs", listDatabases)
	r.HandleFunc("/createtable", createTable)
	r.HandleFunc("/listtables", listTables)
	err = http.ListenAndServe(os.Getenv("URL_HTTP"), r)
	log.Fatalf("server crashed:\n %s\n", err)
}
