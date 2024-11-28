package utils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func SeedAndCleanupData(t *testing.T, seed func(*gorm.DB)) {
	seedData(seed)
	t.Cleanup(cleanupData)
}

func seedData(seed func(*gorm.DB)) {
	db := GetDatabase()
	seed(db)
}

func cleanupData() {
	db := GetDatabase()

	query := "id IS NOT NULL"

	func(models ...interface{}) {
		for _, m := range models {
			db.Where(query).Delete(m)
		}
	}(
		&model.SongSectionType{},
		&model.GuitarTuning{},
		&model.User{},
	)
}
