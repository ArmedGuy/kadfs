package s3

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ArmedGuy/kadfs/kademlia"
	"github.com/gorilla/mux"
)

var kadfs *kademlia.Kademlia

func pathToHash(path string) string {
	hash1 := sha1.New()
	hash1.Write([]byte(path))
	return hex.EncodeToString(hash1.Sum(nil))
}

func ConfigureAndListen(state *kademlia.Kademlia, address string) {
	kadfs = state
	router := mux.NewRouter()

	router.HandleFunc("/kadfs/", listBucket).Methods("GET").Queries("list-type", "2")
	router.HandleFunc("/kadfs/{path:.+$}", getObject).Methods("GET")
	router.HandleFunc("/kadfs/{path:.+$}", putObject).Methods("PUT")
	router.HandleFunc("/kadfs/{path:.+$}", headObject).Methods("HEAD")
	router.HandleFunc("/kadfs/{path:.+$}", deleteObject).Methods("DELETE")

	log.Println("[INFO] s3: Listening, accepting connections at", address)
	log.Fatal(http.ListenAndServe(address, router))
}

func listBucket(w http.ResponseWriter, r *http.Request) {

}

func getObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	hash := pathToHash(path)
	data, _, ok := kadfs.FindValue(hash)
	if !ok {
		w.WriteHeader(404)
	} else {
		w.Write(data)
	}
}

func headObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	hash := pathToHash(path)
	data, _, ok := kadfs.FindValue(hash)
	if !ok {
		w.WriteHeader(404)
	} else {
		w.Header().Add("Content-Length", fmt.Sprintf("%v", len(data)))
		w.Write(data)
	}
}

func putObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	hash := pathToHash(path)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	} else {
		kadfs.Store(hash, data)
		w.WriteHeader(200)
	}
}

func deleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	hash := pathToHash(path)
	kadfs.DeleteValue(hash)
	w.WriteHeader(200)
}
