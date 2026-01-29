package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"microservice-products-catalog/cmd/http/config"
	"microservice-products-catalog/cmd/http/dependencies"
	"microservice-products-catalog/cmd/http/routes"
	"net/http"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()
	dep := dependencies.InitDependencies(cfg)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register your routes
	routes.SetupProductRoutes(mux, dep)
	routes.SetupOrderRoutes(mux, dep)

	const port = ":8000"
	fmt.Printf("Starting server at port %s\n", port)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
