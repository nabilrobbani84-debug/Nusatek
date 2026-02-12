package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"nusatek-backend/internal/config"
	"nusatek-backend/internal/delivery/http"
	"nusatek-backend/internal/repository/postgres"
	redisRepo "nusatek-backend/internal/repository/redis"
	"nusatek-backend/internal/usecase"
	"nusatek-backend/pkg/database"
	"nusatek-backend/pkg/rabbitmq"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect to Database
	db, err := database.ConnectPostgres(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// 3. Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}

	// 4. Connect to RabbitMQ
	rabbitConn, rabbitCh, err := rabbitmq.ConnectRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		// Log but don't fatal, allowing app to run without MQ for demo purposes if needed,
		// though strictly it should probably fail.
		log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
	} else {
		defer rabbitConn.Close()
		defer rabbitCh.Close()
	}

	// 5. Init Layers
	// Timeout context for usecase
	timeoutContext := time.Duration(2) * time.Second

	// Repositories
	propertyRepo := postgres.NewPropertyRepository(db)
	cacheRepo := redisRepo.NewPropertyCacheRepository(rdb)

	// Usecase
	propertyUsecase := usecase.NewPropertyUsecase(propertyRepo, cacheRepo, rabbitCh, timeoutContext)

	// 6. Init Router & Handlers
	r := gin.Default()
	http.NewPropertyHandler(r, propertyUsecase)

	// 5. Run Server
	log.Println("Server running on port", cfg.AppPort)
	if err := r.Run(cfg.AppPort); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
