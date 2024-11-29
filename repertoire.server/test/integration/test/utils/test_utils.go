package utils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"repertoire/server/internal"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	"testing"
)

func GetDatabase() *gorm.DB {
	db, _ := gorm.Open(postgres.Open(core.Dsn))
	return db
}

func GetEnv() internal.Env {
	return internal.NewEnv()
}

func SeedAndCleanupData(t *testing.T, users []model.User, seed func(*gorm.DB)) {
	seedData(seed)
	t.Cleanup(func() {
		cleanupData(users)
	})
}

func seedData(seed func(*gorm.DB)) {
	db := GetDatabase()
	seed(db)
}

func cleanupData(users []model.User) {
	db := GetDatabase()

	for _, user := range users {
		db.Select(clause.Associations).Delete(user)
	}
}
