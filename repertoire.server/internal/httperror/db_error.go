package httperror

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	pgForeignKeyViolation       = "23503"
	pgUniqueViolation           = "23505"
	pgNotNullViolation          = "23502"
	pgCheckViolation            = "23514"
	pgInvalidTextRepresentation = "22P02"
)

type DatabaseErrorOptions struct {
	OnForeignKeyViolation func(error) *ErrorCode // default: NotFoundError
	OnUniqueViolation     func(error) *ErrorCode // default: ConflictError
	OnNotFound            func(error) *ErrorCode // default: NotFoundError
}

// DatabaseError applies sensible defaults: FK violation -> 404,
// unique violation -> 409, missing record -> 404.
func DatabaseError(err error) *ErrorCode {
	return DatabaseErrorWithOptions(err, DatabaseErrorOptions{})
}

// DatabaseErrorWithOptions lets a caller override how specific violation
// types map to ErrorCodes (e.g. a user-FK violation meaning "unauthorized"
// rather than "not found").
func DatabaseErrorWithOptions(err error, opts DatabaseErrorOptions) *ErrorCode {
	onFK := NotFoundError
	onUnique := ConflictError
	onNotFound := NotFoundError
	if opts.OnForeignKeyViolation != nil {
		onFK = opts.OnForeignKeyViolation
	}
	if opts.OnUniqueViolation != nil {
		onUnique = opts.OnUniqueViolation
	}
	if opts.OnNotFound != nil {
		onNotFound = opts.OnNotFound
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return onNotFound(errors.New("record not found"))
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgForeignKeyViolation:
			if column, table, ok := parseForeignKeyDetail(pgErr.Detail); ok {
				msg := fmt.Errorf("%s (%s) does not reference an existing row", table, column)
				if column == "user_id" && table == "users" {
					return ForbiddenError(msg)
				}
				return onFK(msg)
			}
			return onFK(pgMessage(pgErr))
		case pgUniqueViolation:
			return onUnique(pgMessage(pgErr))
		case pgNotNullViolation, pgCheckViolation, pgInvalidTextRepresentation:
			return BadRequestError(pgMessage(pgErr))
		}
	}

	return InternalServerError(err)
}

func pgMessage(pgErr *pgconn.PgError) error {
	if pgErr.Detail != "" {
		return fmt.Errorf("%s: %s", pgErr.Message, pgErr.Detail)
	}
	return errors.New(pgErr.Message)
}

var fkDetailPattern = regexp.MustCompile(`\(([^)]+)\)=\([^)]*\) is not present in table "([^"]+)"`)

func parseForeignKeyDetail(detail string) (column, table string, ok bool) {
	m := fkDetailPattern.FindStringSubmatch(detail)
	if m == nil {
		return "", "", false
	}
	return m[1], m[2], true
}
