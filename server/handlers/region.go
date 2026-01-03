package handlers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"fmt"

	"server/models"
)

func GetRegionsHandler(pgdb *sql.DB, w http.ResponseWriter, r *http.Request) {
	regions, err := models.GetRegionList(pgdb)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(regions); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}