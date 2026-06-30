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
	accessExp, _ := strconv.Atoi(os.Getenv("EXP"))
	refreshExp, _ := strconv.Atoi(os.Getenv("REFRESH_EXP"))

	tokenData := service.TokenInfo{
		Access:     os.Getenv("SECRET"),
		AccessExp:  accessExp,
		Refresh:    os.Getenv("REFRESH"),
		RefreshExp: refreshExp,
	}

	ur := repository.New(pool)
	us := service.New(ur, tokenData)
	uh := handler.NewAuth(us)
	mh := handler.NewMe(us)

	rtr := router.New(uh, mh)

	if err := rtr.Run(os.Getenv("APP_PORT")); err != nil {
		log.Fatal("Server didn't start")
	}

}
