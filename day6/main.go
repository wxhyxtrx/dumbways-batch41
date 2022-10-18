package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World ini halaman home"))
	})
	route.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World ini halaman Project"))
	})
	route.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World ini halaman Contact"))
	})

	fmt.Println("Sever Running on Port 5000")
	http.ListenAndServe("localhost:5000", route)
}
