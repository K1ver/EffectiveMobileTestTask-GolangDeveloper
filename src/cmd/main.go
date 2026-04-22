package main

import (
	"fmt"

	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/config"
	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/handler"
	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/repository"
	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/docs"
)

// @title Subscription Service API
// @version 1.0
// @description REST-сервис для агрегации данных об онлайн подписках
// @host localhost:8080
// @BasePath /
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}
	defer db.Close()
	log.Info("Connected to database")

	// Run migrations
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver: ", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Failed to init migrations: ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations: ", err)
	}
	log.Info("Migrations applied")

	repo := repository.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewHandler(svc)

	r := gin.Default()
	h.RegisterRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
