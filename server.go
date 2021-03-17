package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

////go:embed assets/static
var staticFiles embed.FS

//go:embed assets/static/style.css
var styling []byte

//go:embed assets/templates
var templateFiles embed.FS

type serverCmd struct {
	Port string `help:"listen port" default:":8080"`
}

func (c *serverCmd) Run(ctx *runctx) error {

	indexT, err := template.ParseFS(templateFiles, "assets/templates/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// WTF
	// r.Handle("/x", http.StripPrefix("assets/static", http.FileServer(http.FS(staticFiles))))
	// r.Handle("/x", http.FileServer(http.FS(staticFiles)))
	r.HandleFunc("/s/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Write(styling)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexT.Execute(w, nil)
	})
	r.HandleFunc("/{address}", func(w http.ResponseWriter, r *http.Request) {

	})

	server := handlers.LoggingHandler(os.Stdout, r)

	return http.ListenAndServe(c.Port, server)
}
