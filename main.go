package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		log.Printf("Received post")
	}
	title := r.URL.Path[1:]
	err := templates.ExecuteTemplate(w, "index", title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	//handle static files, perhaps better left to nginx in future?
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//default port, but allow command line argument to override port
	port := ":8080"
	if len(os.Args) > 1 {
		tport := os.Args[1]
		i, err := strconv.Atoi(tport)
		if err != nil || i < 0 || i > 65535 {
			log.Printf("Error: specified port %q is invalid.\n", tport)
			log.Printf("Defaulting to port 8080.\n")
		} else {
			port = ":" + tport
		}
	}
	log.Printf("Listening on port: " + port[1:])

	log.Fatal(http.ListenAndServe(port, nil))
}
