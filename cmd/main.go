package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"wschat/internal/handler"
	"wschat/internal/repository"
	"wschat/internal/router"
	"wschat/internal/service"
	auth_token "wschat/internal/service/auth_token"

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
	accessExp, _ := strconv.Atoi(os.Getenv("ACCESS_EXP"))
	refreshExp, _ := strconv.Atoi(os.Getenv("REFRESH_EXP"))
	accessSecret := os.Getenv("ACCESS")
	refreshSecret := os.Getenv("REFRESH")

	tm := auth_token.NewManager(accessSecret, refreshSecret, accessExp, refreshExp)

	ur := repository.New(pool)
	us := service.New(ur, tm)
	uh := handler.NewAuth(us)
	mh := handler.NewMe(us, tm)

	rtr := router.New(uh, mh)

	if err := rtr.Run(os.Getenv("APP_PORT")); err != nil {
		log.Fatal("Server didn't start")
	}

}
