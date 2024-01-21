package ping

import (
	"net/http"
)

type Pinger interface {
	Ping() error
}

// DBConHandler checks the connection to the database.
func DBConHandler(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
