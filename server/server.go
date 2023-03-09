package main

import (
	"log"
	"net/http"
	"time"
)

func main () {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}


func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request inicializado")
	defer log.Println("Request finalizado")

	select {
	case <-time.After(5 * time.Second):
		log.Println("Request executado com sucesso")
		w.Write([]byte("Request executado com sucesso\n"))

	case <-ctx.Done():
		log.Println("Request cancelado pelo cliente")
		// http.Error(w, "Request cancelado pelo cliente", http.StatusRequestTimeout)
	}
}
