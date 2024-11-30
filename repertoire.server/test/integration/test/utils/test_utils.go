package utils

import (
	"repertoire/server/internal"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetDatabase() *gorm.DB {
	db, _ := gorm.Open(postgres.Open(core.Dsn))
	return db
}

func GetEnv() internal.Env {
	return internal.NewEnv()
}

func CreateValidToken(user model.User) string {
	env := GetEnv()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": env.JwtIssuer,
		"aud": env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	return token
}

func CreateCustomToken(sub string, jti string) string {
	env := GetEnv()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": jti,
		"sub": sub,
		"iss": env.JwtIssuer,
		"aud": env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	return token
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
