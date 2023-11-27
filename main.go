package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcb1310/bca-go-w-test/routes/api"
	"github.com/gorilla/handlers"
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
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	h := api.NewRouter()

	fileServer := http.FileServer(http.Dir("./dist/"))
	h.PathPrefix("/css/").Handler(http.StripPrefix("/css/", fileServer))
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CORS(originsOk, headersOk, methodsOk)(h)))
}
