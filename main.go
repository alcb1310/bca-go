package main

import (
	"log"
	"net/http"

	"github.com/alcb1310/bca-go-w-test/routes"
)

func main() {
	port := ":8000"
	h := routes.NewRouter()

	log.Panic(http.ListenAndServe(port, h))
}
