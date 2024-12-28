package helpers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

var pgErrorToHTTPStatus = map[string]int{
	"23505": http.StatusBadRequest,          // Unique violation
	"23503": http.StatusBadRequest,          // Foreign key violation
	"22001": http.StatusBadRequest,          // String data right truncation
	"23514": http.StatusBadRequest,          // Check violation
	"42P01": http.StatusInternalServerError, // Undefined table
	"22P02": http.StatusBadRequest,          // Invalid text representation
}

func TranslatePgError(err error) (int, string) {
	if err == nil {
		return -1, ""
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if status, found := pgErrorToHTTPStatus[pgErr.Code]; found {
			return status, pgErr.Message
		}
		return http.StatusInternalServerError, "Неизвестная ошибка базы данных"
	}
	return http.StatusInternalServerError, "Ошибка выполнения запроса"
}
