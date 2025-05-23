package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"pvz-service/cmd/app"
	"pvz-service/internal/config"
	"pvz-service/internal/db"
	grpcserver "pvz-service/internal/grpc"
)

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :9000")
		if err := http.ListenAndServe(":9000", nil); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()
}

func startGRPCServerAsync(db *sql.DB, port string) {
	go func() {
		if err := grpcserver.StartGRPCServer(db, port); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := config.LoadConfig()

	database, err := db.InitializeDB(cfg.DbDSN)
	if err != nil {
		log.Fatal("Failed to initialize DB:", err)
	}
	defer database.Close()

	startGRPCServerAsync(database, "3000")

	application := app.MakeApp(database, cfg)

	startMetricsServer()

	log.Printf("Server listening on port %s", cfg.Port)
	log.Fatal(application.Listen(fmt.Sprintf("0.0.0.0:%s", cfg.Port)))
}
