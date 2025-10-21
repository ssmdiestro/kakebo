package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	apihttp "kakebo/internal/http"
	"kakebo/internal/repository"
	"kakebo/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI no definido")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repository.NewMongo(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Error conectando a Mongo: %v", err)
	}

	svc := service.New(db)

	r := apihttp.NewRouter(svc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("API escuchando en :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
