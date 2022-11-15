package router

import (
	"database/sql"
	"net/http"
	"github.com/TechBowl-japan/go-stations/handler"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz/", handler.NewHealthzHandler().ServeHTTP)
	// 下記の2パターンで実行可能
	// mux.HandleFunc("/healthz/", handler.NewHealthzHandler().ServeHTTP)
	// mux.HandleFunc("/healthz/", new(handler.HealthzHandler).ServeHTTP)
	return mux
}
