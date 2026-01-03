package handlers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"fmt"

	"server/models"
)

type Coordinates struct {
    Latitude string `json:"lat"`
    Longitude string `json:"lon"`
}

func GetStopsHandler(pgdb *sql.DB, w http.ResponseWriter, r *http.Request, region string) {
    stops, err := models.GetStopList(pgdb, region)
    if err != nil {
        if err == sql.ErrNoRows {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"error":"no stops found"}`))
            return
        }

        fmt.Println(err)
        http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(stops); err != nil {
        http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
    }
}

func FindNearestStopHandler(pgdb *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }
	
 	var coords Coordinates

	err := json.NewDecoder(r.Body).Decode(&coords)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	stops, err := models.GetStopsAll(pgdb);
	if err != nil {
        if err == sql.ErrNoRows {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"error":"no stops found"}`))
            return
        }

        fmt.Println(err)
        http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
        return
    }

	nearest, err := findNearestStop(coords, stops)
	if err != nil {
        fmt.Println(err)
        http.Error(w, `{"error":"error comparing distances, possibly wrong format of input data"}`, http.StatusInternalServerError)
        return
    }

    fmt.Println(nearest)

	routes, err := models.GetRouteList(pgdb, nearest.Name, nearest.Code)
    if err != nil {
        if err == sql.ErrNoRows {
           w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"stop":    nearest,
				"routes":  []interface{}{},
				"message": "no routes found for this stop",
			})
			return
        }

        fmt.Println(err)
        http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
        return
    }

    sortRoutes(routes)

	response := map[string]interface{}{
        "stop": nearest,
		"routes": routes,
    }

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
