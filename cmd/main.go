package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"wschat/internal/handler"
	"wschat/internal/repository"
	"wschat/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgPath := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"), os.Getenv("DB_NAME"))

	pool, err := pgxpool.New(context.Background(), pgPath)

	if err != nil {
		log.Fatal("Can't connect to DB", err)
	}

	defer pool.Close()
	secret := os.Getenv("SECRET")
	exp, _ := strconv.Atoi(os.Getenv("exp"))

	ur := repository.New(pool)
	us := service.New(ur, secret, exp)
	uh := handler.New(us)

	router := gin.Default()

	handler.RegisterValidators()
	uh.AuthRoutes(router)

	// go func() {
	// 	if err := router.Run(os.Getenv("APP_PORT")); err != nil {
	// 		log.Fatal("Server didn't start")
	// 	}
	// }()

	if err := router.Run(os.Getenv("APP_PORT")); err != nil {
		log.Fatal("Server didn't start")
	}

}
