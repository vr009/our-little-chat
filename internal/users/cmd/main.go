package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	middleware2 "our-little-chatik/internal/middleware"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/proto/users"
	"our-little-chatik/internal/users/internal/delivery"
	"our-little-chatik/internal/users/internal/repo"
	"our-little-chatik/internal/users/internal/usecase"
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

func lookUpDatabaseConfig() *dbConfig {
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

	dbCfg := lookUpDatabaseConfig()

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

	userRepo := repo.NewUserRepo(db)
	useCase := usecase.NewUserUsecase(userRepo)
	userDataHandler := delivery.NewUserEchoHandler(useCase)
	authHandler := delivery.NewAuthEchoHandler(useCase)

	grpcHandler := delivery.NewUserGRPCHandler(useCase)

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

	// Common API
	// Get info about the user which calls the method.
	commonRouter.GET("/me", userDataHandler.GetMe)
	// Deactivate user account which calls the method.
	commonRouter.DELETE("/me", userDataHandler.DeactivateUser)
	// Update user account which calls the method.
	commonRouter.PATCH("/me", userDataHandler.UpdateUser)
	// Search for users using nicknames.
	commonRouter.GET("/search", userDataHandler.SearchUsers)
	// Get user for its ID.
	commonRouter.GET("/:id", userDataHandler.GetUserForID)

	// Auth API
	// Sign up method.
	authRouter.POST("/signup", authHandler.SignUp)
	// Log in method.
	authRouter.POST("/login", authHandler.Login)
	// Log out method.
	authRouter.DELETE("/logout", authHandler.Logout,
		echojwt.WithConfig(config), middleware2.Auth)

	go func() {
		//TODO graceful shutdown + intercepting signals
		usersGRPCPort := os.Getenv("GRPC_USERS_SERVER_PORT")
		if usersGRPCPort == "" {
			panic("no variable GRPC_USERS_SERVER_PORT passed")
		}

		lis, err := net.Listen("tcp", usersGRPCPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		users.RegisterUsersServer(s, grpcHandler)
		log.Printf("grpc server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	e.Logger.Fatal(e.Start(":" + appConfig.Port))
	return nil
}
