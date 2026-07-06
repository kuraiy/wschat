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
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := connectToPostgres()

	if err != nil {
		log.Fatal("Can't connect to DB", err)
	}

	defer pool.Close()

	redisClient := connectToRedis()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Can't connect to redis")
	}

	rdb := repository.NewRedis(redisClient)

	tm := configureTokenManager(*rdb)
	ur := repository.New(pool)
	us := service.New(ur, tm)
	uh := handler.NewAuth(us, tm)
	mh := handler.NewUser(us, tm)

	rtr := router.New(uh, mh, tm)

	if err := rtr.Run(os.Getenv("APP_PORT")); err != nil {
		log.Fatal("Server didn't start")
	}

}

func connectToPostgres() (*pgxpool.Pool, error) {
	pgPath := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"), os.Getenv("DB_NAME"))

	pool, err := pgxpool.New(context.Background(), pgPath)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func connectToRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS")
	redisPass := os.Getenv("REDISPASS")

	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
}

func configureTokenManager(repo repository.Redis) *auth_token.TokenManager {
	accessExp, _ := strconv.Atoi(os.Getenv("ACCESS_EXP"))
	refreshExp, _ := strconv.Atoi(os.Getenv("REFRESH_EXP"))
	accessSecret := os.Getenv("ACCESS")
	refreshSecret := os.Getenv("REFRESH")

	return auth_token.NewManager(accessSecret, refreshSecret, accessExp, refreshExp, repo)
}
