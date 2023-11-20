package main

import (
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	middleware2 "our-little-chatik/internal/middleware"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/proto/users"
	"our-little-chatik/internal/users/internal/delivery"
	"our-little-chatik/internal/users/internal/repo"
	"our-little-chatik/internal/users/internal/usecase"
	"strconv"
	"time"

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

type redisConfig struct {
	Host     string
	Port     string
	Password string
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
	redisCfg := redisConfig{
		Port:     redisPort,
		Host:     redisHost,
		Password: redisPassword,
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host + ":" + redisCfg.Port,
		Password: redisCfg.Password,
	})

	dbCfg := lookUpDatabaseConfig()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	queueURL := os.Getenv("QUEUE_URL")
	if queueURL == "" {
		panic("no QUEUE_URL is passed")
	}
	conn, err := amqp.Dial(queueURL)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("pgx", dbCfg.dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(dbCfg.maxIdleTime)
	db.SetMaxIdleConns(dbCfg.maxIdleConns)
	db.SetMaxOpenConns(dbCfg.maxOpenConns)

	userRepo := repo.NewUserRepo(db)
	sessionRepo := repo.NewSessionRepo(redisClient)
	activationRepo := repo.NewActivationRepo(redisClient)
	mailerRepo := repo.NewMailerQueue(ch)
	useCase := usecase.NewUserUsecase(userRepo, sessionRepo, activationRepo, mailerRepo)
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

	// This middleware check the activated user sessions
	authPlain := middleware2.AuthMiddlewareHandler{
		SessionGetter:       useCase,
		RequiredSessionType: models.PlainSession,
	}

	// This middleware check the non-activated users sessions
	authActivation := middleware2.AuthMiddlewareHandler{
		SessionGetter:       useCase,
		RequiredSessionType: models.ActivationSession,
	}

	// Restricted group
	authRouter := e.Group("/api/v1/auth")
	commonRouter := e.Group("/api/v1/user",
		echojwt.WithConfig(config), authPlain.Auth)

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
		echojwt.WithConfig(config), authPlain.Auth)
	// Activate user after login
	authRouter.POST("/activation", authHandler.Activate,
		echojwt.WithConfig(config), authActivation.Auth)

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
