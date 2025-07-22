package api

import (
	"database/sql"
	"log"
	"net/http"
	"portfolio-api/cmd/service/s3"
)

const (
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := http.NewServeMux()

	//add middleware here later

	// Create a subrouter for /api/v1
	apiV1Router := http.NewServeMux()

	s3Handler := s3.NewHandler()
	s3Handler.RegisterRoutes(apiV1Router)

	// Mount the subrouter at /api/v1/
	router.Handle(APIPrefix+"/", http.StripPrefix(APIPrefix, apiV1Router))

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
