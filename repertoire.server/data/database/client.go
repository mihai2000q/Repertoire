package database

import (
	"fmt"
	"log"
	"repertoire/server/data/logger"
	"repertoire/server/internal"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	*gorm.DB
}

func NewClient(logger *logger.GormLogger, env internal.Env) Client {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		env.DatabasePort,
		env.DatabaseSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: logger,
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	client := Client{db}
	return client
}
