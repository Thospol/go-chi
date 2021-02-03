package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"saaa-api/internal/core/context"
	"saaa-api/internal/core/sql"

	"github.com/go-chi/chi/middleware"
)

// Transaction to do transaction my sql
func Transaction(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		database := sql.Database.Begin()
		context.SetDatabase(r, database)
		next.ServeHTTP(w, r)

		if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
			_ = database.Rollback()

			logEntry := middleware.GetLogEntry(r)
			if logEntry != nil {
				logEntry.Panic(rvr, debug.Stack())
			} else {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()
			}

			w.WriteHeader(http.StatusInternalServerError)
		}

		if context.GetErrMsg(r) != "" || database.Commit().Error != nil {
			_ = database.Rollback()
		}
	}

	return http.HandlerFunc(fn)
}
