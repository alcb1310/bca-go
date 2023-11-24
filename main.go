package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcb1310/bca-go-w-test/routes/api"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	// h := routes.NewRouter()
	h := api.NewRouter()
	handler := c.Handler(h)

	fileServer := http.FileServer(http.Dir("./dist/"))
	h.PathPrefix("/css/").Handler(http.StripPrefix("/css/", fileServer))
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
