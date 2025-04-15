package config

import "time"

type Config struct {
	Server struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	Database struct {
		DSN string
	}
	JWT struct {
		Secret   string
		ExpireAt time.Duration
	}
}

var C *Config

func InitConfig() {
	C = &Config{
		Server: struct {
			Addr         string
			ReadTimeout  time.Duration
			WriteTimeout time.Duration
		}{
			Addr:         ":8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: struct{ DSN string }{
			DSN: "user:pass@tcp(localhost:3306)/cactus?parseTime=true",
		},
		JWT: struct {
			Secret   string
			ExpireAt time.Duration
		}{
			Secret:   "your_jwt_secret_key",
			ExpireAt: 24 * time.Hour,
		},
	}
}
