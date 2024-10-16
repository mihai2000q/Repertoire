package database

import (
	"context"
	"fmt"
	"log"
	"repertoire/models"
	"repertoire/utils"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewClient(lc fx.Lifecycle, env utils.Env) Client {
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
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if env.Environment == utils.DevelopmentEnvironment {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return db.AutoMigrate(
					&models.User{},
					&models.Artist{},
					&models.Playlist{},
					&models.Album{},
					&models.Song{},
				)
			},
		})
	}

	client := Client{
		DB: db,
	}
	return client
}
