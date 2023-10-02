package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	middleware2 "our-little-chatik/internal/middleware"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/user_data/internal/delivery"
	"our-little-chatik/internal/user_data/internal/repo"
	"our-little-chatik/internal/user_data/internal/usecase"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type AppConfig struct {
	Port string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

var (
	defaultMaxOpenConns = 10
	defaultMaxIdleConns = 10
	defaultMaxIdleTime  = time.Minute * 10
)

func getConnectionString() *dbConfig {
	dbCfg := &dbConfig{}
	key, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic(errors.New("connection string not found"))
	}
	dbCfg.dsn = key

	key, ok = os.LookupEnv("DATABASE_MAX_OPEN_CONNS")
	if !ok {
		dbCfg.maxOpenConns = defaultMaxOpenConns
	} else {
		val, err := strconv.Atoi(key)
		if err != nil {
			panic(err.Error())
		}
		dbCfg.maxOpenConns = val
	}

	key, ok = os.LookupEnv("DATABASE_MAX_IDLE_CONNS")
	if !ok {
		dbCfg.maxIdleConns = defaultMaxIdleConns
	} else {
		val, err := strconv.Atoi(key)
		if err != nil {
			panic(err.Error())
		}
		dbCfg.maxIdleConns = val
	}

	key, ok = os.LookupEnv("DATABASE_MAX_IDLE_TIME")
	if !ok {
		dbCfg.maxIdleTime = defaultMaxIdleTime
	} else {
		duration, err := time.ParseDuration(key)
		if err != nil {
			panic(err.Error())
		}
		dbCfg.maxIdleTime = duration
	}

	return dbCfg
}

func main() {
	log.Fatal(run())
}

func run() error {
	appConfig := &AppConfig{}
	port := os.Getenv("USER_DATA_PORT")
	if port == "" {
		panic("empty USER_DATA_PORT")
	}
	appConfig.Port = port

	dbCfg := getConnectionString()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	pool, err := pgxpool.New(context.Background(), dbCfg.dsn)
	if err != nil {
		log.Fatal("ERROR: : " + err.Error())
	} else {
		slog.Info("Connected to postgres: %s", dbCfg.dsn)
	}
	defer pool.Close()

	db, err := sql.Open("pgx", dbCfg.dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(dbCfg.maxIdleTime)
	db.SetMaxIdleConns(dbCfg.maxIdleConns)
	db.SetMaxOpenConns(dbCfg.maxOpenConns)

	personRepo := repo.NewPersonRepo(db)
	useCase := usecase.NewUserUsecase(personRepo)
	userDataHandler := delivery.NewUserEchoHandler(useCase)
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
	adminRouter := e.Group("/api/v1/admin")

	// Admin CRUD API
	adminRouter.PATCH("/user", userDataHandler.UpdateUser, middleware2.AdminAuth)
	adminRouter.DELETE("/user", userDataHandler.DeleteUser, middleware2.AdminAuth)
	adminRouter.GET("/user", userDataHandler.GetUser, middleware2.AdminAuth)

	// Common API
	commonRouter.GET("/all", userDataHandler.GetAllUsers)
	commonRouter.GET("/me", userDataHandler.GetMe)
	commonRouter.GET("/search", userDataHandler.FindUser)
	commonRouter.DELETE("/:id/info", userDataHandler.DeleteUser)
	commonRouter.PATCH("/info", userDataHandler.UpdateUser)
	commonRouter.GET("/:id/info", userDataHandler.GetUser)

	// Auth API
	authRouter.POST("/signup", authHandler.SignUp)
	authRouter.POST("/login", authHandler.Login)
	authRouter.DELETE("/logout", authHandler.Logout,
		echojwt.WithConfig(config), middleware2.Auth)

	e.Logger.Fatal(e.Start(":" + appConfig.Port))
	return nil
}
