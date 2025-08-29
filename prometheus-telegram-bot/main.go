package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/alert", alertHandler)

	log.Println("🚀 Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
