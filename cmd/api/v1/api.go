package api

import (
	"net/http"
	Logger "velocityApi/logs"
	"velocityApi/services/velocity"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type APIServer struct {
	addr           string
	clickhouseConn driver.Conn
	logger         *Logger.ApiLogger
}

func NewAPIServer(addr string, db driver.Conn, logger *Logger.ApiLogger) *APIServer {
	return &APIServer{
		addr:           addr,
		clickhouseConn: db,
		logger:         logger,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	dataHandler := velocity.NewHandler(&s.clickhouseConn, s.logger)
	dataHandler.RegisterRoutes(router)
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", router))
	s.logger.Write("Listening on:" + s.addr)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	return server.ListenAndServe()
}
