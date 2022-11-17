package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz/", handler.NewHealthzHandler().ServeHTTP)
	// 下記の2パターンで実行可能
	// mux.HandleFunc("/healthz/", handler.NewHealthzHandler().ServeHTTP)
	// mux.HandleFunc("/healthz/", new(handler.HealthzHandler).ServeHTTP)

	mux.HandleFunc("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)).ServeHTTP)
	return mux
}
