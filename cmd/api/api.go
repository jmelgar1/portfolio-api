package api

import (
	"database/sql"
	"log"
	"net/http"
	"portfolio-api/cmd/service/user"
)

const (
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	//add middleware here later

	// Create a subrouter for /api/v1
	apiV1Router := http.NewServeMux()
	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(apiV1Router)

	// Mount the subrouter at /api/v1/
	router.Handle(APIPrefix, http.StripPrefix(APIPrefix, apiV1Router))

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
