package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
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
		log.Println("no .env file, reading from environment")
	}

	pool, err := connectToPostgres()

	if err != nil {
		log.Fatal("Can't connect to DB", err)
	}

	defer pool.Close()

	redisClient, err := connectToRedis()

	if err != nil {
		log.Fatal("Can't connect to Redis", err)
	}

	rdb := repository.NewRedis(redisClient)

	tm := configureTokenManager(rdb)
	ur := repository.New(pool)
	us := service.New(ur, tm)
	uh := handler.NewAuth(us, tm)
	mh := handler.NewUser(us, tm)

	rtr := router.New(uh, mh, tm)

	if err := rtr.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
		log.Fatal("Server didn't start")
	}

}

func connectToPostgres() (*pgxpool.Pool, error) {
	pgPath := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"),
		os.Getenv("PGPORT"), os.Getenv("DB_NAME"))

	pool, err := pgxpool.New(context.Background(), pgPath)

	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func connectToRedis() (*redis.Client, error) {
	redisAddr := os.Getenv("REDIS")
	redisPass := os.Getenv("REDISPASS")

	rClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	if err := rClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return rClient, nil
}

func configureTokenManager(repo *repository.Redis) *auth_token.TokenManager {
	accessExp, _ := strconv.Atoi(os.Getenv("ACCESS_EXP"))
	refreshExp, _ := strconv.Atoi(os.Getenv("REFRESH_EXP"))
	accessSecret := os.Getenv("ACCESS")
	refreshSecret := os.Getenv("REFRESH")

	return auth_token.NewManager(accessSecret, refreshSecret, accessExp, refreshExp, repo)
}
