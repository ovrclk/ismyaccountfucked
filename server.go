package main

import (
	"embed"
	"fmt"
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

	showT, err := template.ParseFS(templateFiles, "assets/templates/show.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// WTF
	// r.Handle("/x", http.StripPrefix("assets/static", http.FileServer(http.FS(staticFiles))))
	// r.Handle("/x", http.FileServer(http.FS(staticFiles)))
	r.HandleFunc("/s/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(styling)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexT.Execute(w, nil)
	})

	r.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		address := r.FormValue("address")
		if address == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/a/"+address, http.StatusSeeOther)
		return
	}).Methods("POST")

	r.HandleFunc("/a/{address}", func(w http.ResponseWriter, r *http.Request) {
		address, ok := mux.Vars(r)["address"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		status, err := getStatus(ctx.cctx, ctx.denom, address)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		if err := showT.Execute(w, status); err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

	})

	server := handlers.LoggingHandler(os.Stdout, r)

	fmt.Printf("running server on port %v\n\n", c.Port)

	return http.ListenAndServe(c.Port, server)
}
