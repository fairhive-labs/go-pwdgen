package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request from", r.RemoteAddr)
		fmt.Fprint(w, "TRENDEV on Heroku!")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
