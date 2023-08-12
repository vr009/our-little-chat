package main

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	middleware2 "our-little-chatik/internal/middleware"
	"our-little-chatik/internal/pkg"
	"strconv"

	"our-little-chatik/internal/user_data/internal/delivery"
	"our-little-chatik/internal/user_data/internal/repo"
	"our-little-chatik/internal/user_data/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type AppConfig struct {
	Port int
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
	configPath := os.Getenv("USER_DATA_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("user-data-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file", configPath)
	}

	appConfig := &AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := GetConnectionString()
	if err != nil {
		log.Fatal(err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("ERROR: : " + err.Error())
	} else {
		slog.Info("Connected to postgres: %s", connString)
	}
	defer pool.Close()

	personRepo := repo.NewPersonRepo(pool)
	useCase := usecase.NewUserdataUseCase(personRepo)
	userDataHandler := delivery.NewUserdataEchoHandler(useCase)
	authHandler := delivery.NewAuthEchoHandler(useCase)

	e := echo.New()
	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())

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

	// Restricted group
	authRouter := e.Group("/api/v1/auth")
	commonRouter := e.Group("/api/v1/user", echojwt.WithConfig(config), middleware2.Auth)
	adminRouter := e.Group("/api/v1/admin", echojwt.WithConfig(config), middleware2.Auth)

	// CRUD API
	adminRouter.POST("/user/new", userDataHandler.CreateUser)
	adminRouter.POST("/user", userDataHandler.UpdateUser)
	adminRouter.DELETE("/user", userDataHandler.DeleteUser)
	adminRouter.GET("/user", userDataHandler.GetUser)

	// Common API
	commonRouter.GET("/all", userDataHandler.GetAllUsers)
	commonRouter.GET("/me", userDataHandler.GetMe)
	commonRouter.GET("/search", userDataHandler.FindUser)

	// Auth API
	authRouter.POST("/signup", authHandler.SignUp)
	authRouter.POST("/login", authHandler.Login)
	authRouter.DELETE("/logout", authHandler.Logout,
		echojwt.WithConfig(config), middleware2.Auth)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(appConfig.Port)))
	return nil
}
