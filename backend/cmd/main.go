package main

import (
	"adv/internal/config"
	db2 "adv/internal/db"
	"adv/internal/handlers"
	"adv/internal/models"
	"adv/internal/repository"
	"adv/internal/router"
	"adv/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}
func getEnvInt(key string, fallback int) int {
	if valueStr, ok := os.LookupEnv(key); ok {
		if valueInt, err := strconv.Atoi(valueStr); err == nil {
			return valueInt
		}
	}
	return fallback
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	logger.InitLogger()
	appConfig := config.Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DbHost:        getEnv("DB_HOST", "localhost"),
		DbUser:        getEnv("DB_USER", "root"),
		DbPassword:    getEnv("DB_PASSWORD", ""),
		DbName:        getEnv("DB_NAME", "your_db"),
		DbPort:        getEnv("DB_PORT", "5432"),
		DbSSLMode:     getEnv("DB_SSLMODE", "disable"),
		RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
		Redispassword: getEnv("REDIS_PASSWORD", ""),
		RedisDb:       getEnvInt("REDIS_DB", 0),
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		SMTPHost:      getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
	}

	db, err := db2.ConnectDB(&appConfig)
	if err != nil {
		logger.GetLogger().Fatal("Failed to connect to database:", err)
	}
	if err = db.AutoMigrate(&models.User{}, &models.Book{}, &models.Subscribers{}, &models.UserBook{}); err != nil {
		logger.GetLogger().Fatal("Failed to migrate database:", err)
	}

	repo := repository.NewRepository(db)
	userhandler := handlers.NewUserHandlers(repo, appConfig)
	bookhandler := handlers.NewBookHandlers(repo, appConfig)

	r := gin.Default()
	router := router.NewRourer(*userhandler, *bookhandler)
	router.Setup(r)

	srv := &http.Server{
		Addr:    ":" + appConfig.AppPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("ListenAndServe:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.GetLogger().Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.GetLogger().Fatal("Server forced to shutdown:", err)
	}

	logger.GetLogger().Println("Server exiting")
}
