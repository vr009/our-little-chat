package models

type DBConfig struct {
	Port     string
	Password string
	DB       int
	TtlHours int
}

type AppConfig struct {
	Address  string
	DataBase DBConfig
}
