package main

import (
	"fmt"
	"log"
	"database/sql"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"server/database"
	"server/handlers"
)

var pgdb *sql.DB
var r *mux.Router

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}
	pgdb = database.SetDB(pgdb)
	//database.ApplyMigrations()
	defer pgdb.Close()

	r = mux.NewRouter()

	// endpoints
	r.HandleFunc("/regions", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Accessed GET /regions\n")
		handlers.GetRegionsHandler(pgdb, w, r)
	}).Methods("GET")

	r.HandleFunc("/stops", func(w http.ResponseWriter, r *http.Request) {
    region := r.URL.Query().Get("region")
    log.Printf("Accessed GET /stops?region=%s\n", region)
		handlers.GetStopsHandler(pgdb, w, r, region)
	}).Methods("GET")

	r.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
    stop := r.URL.Query().Get("stop")
	code := r.URL.Query().Get("code")
    log.Printf("Accessed GET /routes?stop=%s&code=%s\n", stop, code)
		handlers.GetRoutesHandler(pgdb, w, r, stop, code)
	}).Methods("GET")

	r.HandleFunc("/arrivals", func(w http.ResponseWriter, r *http.Request) {
    stop := r.URL.Query().Get("stop")
	code := r.URL.Query().Get("code")
	route := r.URL.Query().Get("route")
    log.Printf("Accessed GET /arrivals?stop=%s&code=%s&route=%s\n", stop, code, route)
		handlers.GetArrivalsHandler(pgdb, w, r, stop, code, route)
	}).Methods("GET")

	r.HandleFunc("/nearest_stop", func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Accessed POST /nearest_stop\n")
		handlers.FindNearestStopHandler(pgdb, w, r)
	}).Methods("POST")

	// catch all OPTIONS requests, so CORS middleware is executed
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Use(corsMiddleware)

	server := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Println("\nStarting server at http://127.0.0.1:8080/")
	fmt.Printf("Quit the server with CONTROL-C.\n\n")

	log.Fatal(server.ListenAndServe())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originsEnv := os.Getenv("ALLOWED_CLIENT_ORIGINS")
		if originsEnv == "" {
			log.Fatal("ALLOWED_CLIENT_ORIGINS not set in environment")
		}

		// Split by comma into slice
		allowedOrigins := strings.Split(originsEnv, ",")

		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}

		// Get the request origin
		origin := r.Header.Get("Origin")

		// Allow the origin only if it matches one of the allowed ones
		if isAllowedOrigin(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == origin {
			return true
		}
	}
	return false
}