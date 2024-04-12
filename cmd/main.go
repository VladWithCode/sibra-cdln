package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/vladwithcode/sibra-cdln/internal/routes"
)

func main() {
	godotenv.Load(".env")

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	addr := fmt.Sprintf(":%v", portStr)
	fmt.Printf("Server listening on http://localhost%v\n", addr)

	router := routes.NewRouter()

	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}
