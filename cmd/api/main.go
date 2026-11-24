package main

import (
	"log"
	"net/http"

	"github.com/MattP2001/envios-demoras-api/internal/handlers"
	"github.com/MattP2001/envios-demoras-api/internal/store"
)

func main() {
	s, err := store.New("envios.db")
	if err != nil {
		log.Fatal(err)
	}
	h := &handlers.Handler{Store: s}
	log.Println("escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Routes()))
}
