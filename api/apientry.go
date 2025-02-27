package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StoreData(w http.ResponseWriter, r *http.Request) {
	// zv-webhook will only send post request
	if r.Method != http.MethodPost {
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	process(w, r)
}

// StartHTTPServer starts the HTTP server to listen for data on a specified port
func StartHTTPServer(port string) {
	http.HandleFunc("/zv-webhook", StoreData)

	fmt.Println("HTTP Server started. Listening on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func StartAPI() {
	port := os.Getenv("API_PORT")
	if len(port) < 1 {
		port = "80"
	}

	go StartHTTPServer(port)
}
