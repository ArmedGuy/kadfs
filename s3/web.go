package s3

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ConfigureAndListen(address string) {
	router := mux.NewRouter()
	log.Println("[INFO] s3: Listening, accepting connections at", address)
	log.Fatal(http.ListenAndServe(address, router))
}
