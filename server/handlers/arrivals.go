package handlers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"fmt"

	"server/models"
)

const limit = 5

func GetArrivalsHandler(pgdb *sql.DB, w http.ResponseWriter, r *http.Request, stop string, code string, route string) {
    arrivals, err := models.GetArrivalList(pgdb, stop, code, route)
    if err != nil || len(arrivals) == 0 {
        if err == sql.ErrNoRows {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"error":"no arrivals found"}`))
            return
        }

        fmt.Println(err)
        http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
        return
    }

	nextArrivals := findNextArrivals(arrivals, limit)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(nextArrivals); err != nil {
        http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
    }
}
