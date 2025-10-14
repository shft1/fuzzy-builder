package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

var openapiSpec = `{
  "openapi": "3.0.0",
  "info": { "title": "Fuzzy Builder API", "version": "1.0.0" },
  "paths": {}
}`

func (s *Server) registerSwagger(r *mux.Router) {
	r.HandleFunc("/swagger/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(openapiSpec))
	}).Methods(http.MethodGet)
	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><body><h1>API</h1><a href="/swagger/openapi.json">openapi.json</a></body></html>`))
	}).Methods(http.MethodGet)
}
