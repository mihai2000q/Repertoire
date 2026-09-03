package httperror

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"repertoire/server/internal/httperror"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	fkViolationCode      = "23503"
	uniqueViolationCode  = "23505"
	notNullViolationCode = "23502"
	checkViolationCode   = "23514"
	invalidTextCode      = "22P02"
	unmappedCode         = "40001" // serialization_failure — deliberately not handled
)

func TestDatabaseError_WhenErrorIsRecordNotFound_ShouldReturnNotFoundError(t *testing.T) {
	errCode := httperror.DatabaseError(gorm.ErrRecordNotFound)

	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.EqualError(t, errCode.Error, "record not found")
}

func TestDatabaseError_WhenErrorWrapsRecordNotFound_ShouldReturnNotFoundError(t *testing.T) {
	wrapped := fmt.Errorf("get by id: %w", gorm.ErrRecordNotFound)

	errCode := httperror.DatabaseError(wrapped)

	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.EqualError(t, errCode.Error, "record not found")
}

func TestDatabaseError_WhenErrorIsForeignKeyViolation_ShouldReturnNotFoundError(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    fkViolationCode,
		Message: `insert or update on table "song_sections" violates foreign key constraint "fk_song_sections_song"`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.EqualError(t, errCode.Error, pgErr.Message)
}

func TestDatabaseError_WhenErrorIsUniqueViolation_ShouldReturnConflictError(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    uniqueViolationCode,
		Message: `duplicate key value violates unique constraint "uq_users_email"`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.EqualError(t, errCode.Error, pgErr.Message)
}

func TestDatabaseError_WhenErrorIsNotNullViolation_ShouldReturnBadRequestError(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    notNullViolationCode,
		Message: `null value in column "name" of relation "song_section_types" violates not-null constraint`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.EqualError(t, errCode.Error, pgErr.Message)
}

func TestDatabaseError_WhenErrorIsCheckViolation_ShouldReturnBadRequestError(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    checkViolationCode,
		Message: `new row for relation "songs" violates check constraint "chk_songs_difficulty"`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.EqualError(t, errCode.Error, pgErr.Message)
}

func TestDatabaseError_WhenErrorIsInvalidTextRepresentation_ShouldReturnBadRequestError(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    invalidTextCode,
		Message: `invalid input syntax for type uuid: "not-a-uuid"`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.EqualError(t, errCode.Error, pgErr.Message)
}

func TestDatabaseError_WhenErrorIsUnmappedPgError_ShouldReturnInternalServerError(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    unmappedCode,
		Message: "could not serialize access due to concurrent update",
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.EqualError(t, errCode.Error, pgErr.Error()) // passed through untouched
}

func TestDatabaseError_WhenErrorIsGeneric_ShouldReturnInternalServerError(t *testing.T) {
	genericErr := errors.New("connection refused")

	errCode := httperror.DatabaseError(genericErr)

	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.EqualError(t, errCode.Error, "connection refused")
}

func TestDatabaseError_WhenPgErrorHasDetail_ShouldJoinMessageAndDetail(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    fkViolationCode,
		Message: "insert or update violates foreign key constraint",
		Detail:  `Key (song_id)=(3fa85f64-5717-4562-b3fc-2c963f66afa6) is not present in table "songs".`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.EqualError(t, errCode.Error, pgErr.Message+": "+pgErr.Detail)
}

func TestDatabaseError_WhenPgErrorHasNoDetail_ShouldUseMessageOnly(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    fkViolationCode,
		Message: `insert or update on table "songs" violates foreign key constraint "fk_songs_album"`,
	}

	errCode := httperror.DatabaseError(pgErr)

	assert.EqualError(t, errCode.Error, pgErr.Message)
}

func TestDatabaseErrorWithOptions_WhenForeignKeyViolationAndCustomHandlerProvided_ShouldUseCustomHandler(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    fkViolationCode,
		Message: `insert or update on table "song_section_types" violates foreign key constraint "fk_song_section_types_user"`,
	}

	// Mirrors the real use case: a FK violation on the caller's own user ID
	// reads as "your session is invalid", not "resource not found".
	errCode := httperror.DatabaseErrorWithOptions(pgErr, httperror.DatabaseErrorOptions{
		OnForeignKeyViolation: httperror.UnauthorizedError,
	})

	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
}

func TestDatabaseErrorWithOptions_WhenUniqueViolationAndCustomHandlerProvided_ShouldUseCustomHandler(t *testing.T) {
	pgErr := &pgconn.PgError{
		Code:    uniqueViolationCode,
		Message: `duplicate key value violates unique constraint "uq_playlists_share_slug"`,
	}

	errCode := httperror.DatabaseErrorWithOptions(pgErr, httperror.DatabaseErrorOptions{
		OnUniqueViolation: httperror.ForbiddenError,
	})

	assert.Equal(t, http.StatusForbidden, errCode.Code)
}

func TestDatabaseErrorWithOptions_WhenRecordNotFoundAndCustomHandlerProvided_ShouldUseCustomHandler(t *testing.T) {
	// e.g. hiding whether a resource exists at all from an unauthorized caller.
	errCode := httperror.DatabaseErrorWithOptions(gorm.ErrRecordNotFound, httperror.DatabaseErrorOptions{
		OnNotFound: httperror.ForbiddenError,
	})

	assert.Equal(t, http.StatusForbidden, errCode.Code)
}

func TestDatabaseErrorWithOptions_WhenOptionsAreEmpty_ShouldFallBackToDefaults(t *testing.T) {
	pgErr := &pgconn.PgError{Code: uniqueViolationCode, Message: "duplicate key"}

	errCode := httperror.DatabaseErrorWithOptions(pgErr, httperror.DatabaseErrorOptions{})

	assert.Equal(t, http.StatusConflict, errCode.Code)
}
