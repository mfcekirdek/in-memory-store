// Package internal has private application and library code.
package internal

import (
	"fmt"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/handler"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/middleware"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/repository"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/service"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/utils"
	"net/http"

	"gitlab.com/mfcekirdek/in-memory-store/configs"
)

// Server contains mux and config objects.
type Server struct {
	mux    *http.ServeMux
	config *configs.Config
}

// Takes config parameters and creates a new Server instance with given configs.
func NewServer(c *configs.Config) *Server {
	server := &Server{
		mux:    http.NewServeMux(),
		config: c,
	}
	return server
}

// Creates all the routings.
// Wraps handler instance with Header and Logger middlewares.
// Starts the server with given configs.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Server.Port)
	s.Routes()
	var wrappedMux http.Handler
	wrappedMux = middleware.NewHeaderMiddleware(s.mux)
	if s.config.IsDebug {
		wrappedMux = middleware.NewLoggerMiddleware(wrappedMux)
	}
	return http.ListenAndServe(addr, wrappedMux)
}

// Creates routings.
func (s *Server) Routes() {
	storeRepository := repository.NewStoreRepository()
	storeService := service.NewStoreService(storeRepository, s.config.SaveToFileInterval, s.config.StorageDirPath)
	storeHandler := handler.NewStoreHandler(storeService)
	s.mux.HandleFunc("/health", checkHealth)
	s.mux.HandleFunc("/api/v1/store", storeHandler.Flush)
	s.mux.Handle("/api/v1/store/", storeHandler)
}

// Responds GET requests to /health with {"status": "OK"} if the application is up and running.
func checkHealth(w http.ResponseWriter, r *http.Request) {
	utils.ReturnJSONResponse(w, r, map[string]string{"status": "OK"}, "healthy")
}
