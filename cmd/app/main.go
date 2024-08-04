package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/a-h/respond"
	"github.com/a-h/rqlite-test/db"
	_ "github.com/rqlite/gorqlite/stdlib"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info("Performing DB migrations")

	databaseURL := db.URL{
		User:     "admin",
		Password: "secret",
		Host:     "localhost",
		Port:     4001,
		Secure:   false,
	}
	driver, err := sql.Open("rqlite", databaseURL.DataSourceName())
	if err != nil {
		log.Error("failed to open database", slog.Any("error", err))
		os.Exit(1)
	}

	if err := db.Migrate(databaseURL); err != nil {
		log.Error("migrations failed", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("Migrations complete")

	log.Info("Starting server", slog.Int("port", 8080))

	queries := db.New(driver)

	mux := http.NewServeMux()
	mux.HandleFunc("/documents", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			docs, err := queries.DocumentsSelectMany(r.Context())
			if err != nil {
				log.Error("failed to list documents", slog.Any("error", err), slog.String("type", reflect.TypeOf(err).String()))
				respond.WithError(w, "failed to list documents", http.StatusInternalServerError)
				return
			}
			respond.WithJSON(w, docs, http.StatusOK)
			return
		}
		if r.Method == http.MethodPost {
			var req DocumentsPostRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				respond.WithError(w, "invalid request", http.StatusBadRequest)
				return
			}
			//TODO: Validate the request.
			if err := queries.DocumentsInsert(r.Context(), db.DocumentsInsertParams{
				Name:    req.Name,
				Content: req.Content,
			}); err != nil {
				log.Error("failed to insert document", slog.Any("error", err))
				respond.WithError(w, "failed to insert document", http.StatusInternalServerError)
				return
			}
			respond.WithJSON(w, "document inserted", http.StatusCreated)
			return
		}
		respond.WithError(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/document/{id}", func(w http.ResponseWriter, r *http.Request) {
		log.Info("GET /document/{id}", slog.String("id", r.PathValue("id")))
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			respond.WithError(w, "invalid id", http.StatusBadRequest)
			return
		}
		doc, err := queries.DocumentsSelectOneByID(r.Context(), id)
		if err != nil {
			log.Error("failed to get document", slog.Any("error", err))
			respond.WithError(w, "failed to get document", http.StatusInternalServerError)
			return
		}
		respond.WithJSON(w, doc, http.StatusOK)
		return
	})
	http.ListenAndServe(":8080", mux)
}

type DocumentsPostRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
