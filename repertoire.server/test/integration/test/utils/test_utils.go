package utils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mime/multipart"
	"os"
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

func AttachFileToMultipartBody(fileName string, formName string, multiWriter *multipart.Writer) {
	tempFile, _ := os.CreateTemp("", fileName)
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())

	fileWriter, _ := multiWriter.CreateFormFile(formName, tempFile.Name())

	file, _ := os.Open(tempFile.Name())
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, _ = file.WriteTo(fileWriter)
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
