package internals

import (
	"fmt"
	"net/http"

	"gitlab.com/mfcekirdek/in-memory-store/internals/handler"
	"gitlab.com/mfcekirdek/in-memory-store/internals/repository"
	"gitlab.com/mfcekirdek/in-memory-store/internals/service"

	"gitlab.com/mfcekirdek/in-memory-store/configs"
	"gitlab.com/mfcekirdek/in-memory-store/internals/middleware"
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
)

type Server struct {
	mux    *http.ServeMux
	config *configs.Config
}

func NewServer(c *configs.Config) *Server {
	server := &Server{
		mux:    http.NewServeMux(),
		config: c,
	}
	return server
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Server.Port)
	s.Routes()
	var wrappedMux http.Handler
	wrappedMux = middleware.NewHeaderMiddleware(s.mux)
	if s.config.IsDebug {
		wrappedMux = middleware.NewLoggerMiddleware(wrappedMux)
	}
	return http.ListenAndServe(addr, wrappedMux)
	//return http.ListenAndServeTLS(addr, "configs/tls/server.crt", "configs/tls/server.key", wrappedMux)
}

func (s *Server) Routes() {
	storeRepository := repository.NewStoreRepository()
	storeService := service.NewStoreService(storeRepository, s.config.SaveToFileInterval, s.config.StorageDirPath)
	storeHandler := handler.NewStoreHandler(storeService)
	s.mux.HandleFunc("/health", checkHealth)
	s.mux.HandleFunc("/api/v1/store", storeHandler.Flush)
	s.mux.Handle("/api/v1/store/", storeHandler)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	utils.ReturnJSONResponse(w, r, map[string]string{"status": "OK"}, "healthy")
}
