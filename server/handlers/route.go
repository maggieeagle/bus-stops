package handlers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"fmt"

	"server/models"
)

func GetRoutesHandler(pgdb *sql.DB, w http.ResponseWriter, r *http.Request, stop string, code string) {
    /* stop, err := strconv.Atoi(stopStr)
	if err != nil {
		http.Error(w, `{"error":"invalid stop ID"}`, http.StatusBadRequest)
		return
	} */

    routes, err := models.GetRouteList(pgdb, stop, code)
    if err != nil || len(routes) == 0 {
        if err == sql.ErrNoRows {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"error":"no routes found"}`))
            return
        }

        fmt.Println(err)
        http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
        return
    }

    sortRoutes(routes)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(routes); err != nil {
        http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
    }
}
