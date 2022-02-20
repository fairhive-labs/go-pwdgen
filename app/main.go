package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request from", r.RemoteAddr)
		fmt.Fprint(w, "TRENDEV x Heroku")
	})
	log.Println("👉 Starting server on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
