package main

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	middleware2 "our-little-chatik/internal/middleware"
	"our-little-chatik/internal/pkg"
	"strconv"

	"our-little-chatik/internal/chat/internal/delivery"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/chat/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type AppConfig struct {
	Port  int
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func GetConnectionString() (string, error) {
	key, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func main() {
	log.Fatal(run())
}

func run() error {
	var err error
	appConfig := AppConfig{}
	appConfig.Port, err = strconv.Atoi(os.Getenv("CHAT_PORT"))
	if err != nil {
		panic(err.Error())
	}
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		panic("empty redis host")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		panic("empty redis port")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		panic("empty redis password")
	}

	appConfig.Redis.Port = redisPort
	appConfig.Redis.Host = redisHost
	appConfig.Redis.Password = redisPassword

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.Redis.Host + ":" + appConfig.Redis.Port,
		Password: appConfig.Redis.Password,
	})

	ctx := context.Background()
	connStr, err := GetConnectionString()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	repop := repo.NewPostgresRepo(pool)
	repoRed := repo.NewRedisRepo(redisClient)
	uc := usecase.NewChatUseCase(repop, repoRed)

	handler := delivery.NewChatEchoHandler(uc)

	e := echo.New()
	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())

	// Restricted group
	r := e.Group("/api/v1")

	key, err := pkg.GetSignedKey()
	if err != nil {
		panic(err.Error())
	}

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(pkg.JwtCustomClaims)
		},
		SigningKey:  []byte(key),
		TokenLookup: "cookie:Token",
	}
	r.Use(echojwt.WithConfig(config), middleware2.Auth)

	// Admin API
	adminRouter := e.Group("/api/v1/admin")
	// Getting chat info
	adminRouter.GET("/chat", handler.GetChat)
	adminRouter.DELETE("/chat", handler.DeleteChat)
	adminRouter.POST("/chat", handler.PostNewChat)

	// Get chat info
	r.GET("/chat", handler.GetChat)
	// Get chat messages
	r.GET("/conv", handler.GetChatMessages)
	// Get the list of users chats
	r.GET("/list", handler.GetChatList)
	// Create a new chat
	r.POST("/new", handler.PostNewChat)
	// Update photo url of the chat
	r.POST("/chat/photo", handler.ChangeChatPhoto)
	// Add users to chat
	r.POST("/chat/users", handler.AddUsersToChat)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(appConfig.Port)))
	return nil
}
