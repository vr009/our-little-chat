package main

import (
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	middleware2 "our-little-chatik/internal/middleware"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/proto/session"
	"our-little-chatik/internal/pkg/proto/users"
	"strconv"
	"time"

	"our-little-chatik/internal/chat/internal/delivery"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/chat/internal/usecase"

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

	dbCfg := lookUpDatabaseConfig()
	db, err := sql.Open("pgx", dbCfg.dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetConnMaxIdleTime(dbCfg.maxIdleTime)
	db.SetMaxIdleConns(dbCfg.maxIdleConns)
	db.SetMaxOpenConns(dbCfg.maxOpenConns)

	usersGRPCHost := os.Getenv("GRPC_USERS_SERVER_HOST")
	if usersGRPCHost == "" {
		panic("no variable GRPC_USERS_SERVER_HOST passed")
	}
	usersGRPCPort := os.Getenv("GRPC_USERS_SERVER_PORT")
	if usersGRPCPort == "" {
		panic("no variable GRPC_USERS_SERVER_PORT passed")
	}

	usersGRPCServerAddr := usersGRPCHost + usersGRPCPort
	// Set up a connection to the server.
	conn, err := grpc.Dial(usersGRPCServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	usersGRPCCl := users.NewUsersClient(conn)

	usersClient := repo.NewUserDataClient(usersGRPCCl)
	repoPostgres := repo.NewPostgresRepo(db)
	repoRedis := repo.NewRedisRepo(redisClient)
	uc := usecase.NewChatUseCase(repoPostgres, repoRedis, usersClient)
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

	auth := middleware2.AuthMiddlewareHandler{
		SessionGetter:       middleware2.NewDefaultGRPCSessionGetter(session.NewSessionClient(conn)),
		RequiredSessionType: models.PlainSession,
	}
	r.Use(echojwt.WithConfig(config), auth.Auth)

	chatRouter := r.Group("/chat")

	// Get chat info
	chatRouter.GET("/:id", handler.GetChat)
	// Get chat messages
	chatRouter.GET("/:id/messages", handler.GetChatMessages)
	// Get the list of users chats
	chatRouter.GET("/list", handler.GetChatList)
	// Create a new chat
	chatRouter.POST("/new", handler.PostNewChat)
	// Update photo url of the chat
	chatRouter.POST("/photo", handler.ChangeChatPhoto)
	// Add users to chat
	chatRouter.POST("/users", handler.AddUsersToChat)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(appConfig.Port)))
	return nil
}
