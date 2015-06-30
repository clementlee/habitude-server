package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL.Query())
	log.Println(r.Method)
	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))
}

type UserReq struct {
	Email string `json:"email"`
}

/**
 * Only handles POST.
 */
func getuser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported", http.StatusMethodNotAllowed)
	}
	decoder := json.NewDecoder(r.Body)
	var t UserReq
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Hello, user at %s\n", t.Email)

}

func cron() {
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		//fmt.Println("hi")
		//run periodic update code.
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/getuser", getuser)

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
	go cron()

	log.Fatal(http.ListenAndServe(port, nil))
}
