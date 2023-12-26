package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleSlackRequest)
	log.Println("Starting Slack bot...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
