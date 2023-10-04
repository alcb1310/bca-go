package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcb1310/bca-go-w-test/routes"
	"github.com/joho/godotenv"
)

func main() {
	port, portRead := os.LookupEnv("PORT")
	if !portRead {
		godotenv.Load()
		port, portRead = os.LookupEnv("PORT")
		if !portRead {
			log.Panic("Unable to load environment variables")
		}
	}

	h := routes.NewRouter()

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", port), h))
}
