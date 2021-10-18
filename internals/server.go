package internals

import (
	"fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"
	"gitlab.com/mfcekirdek/in-memory-store/configs"
)

type Server struct {
	e      *echo.Echo
	config *configs.Config
}

func NewServer(c *configs.Config) *Server {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	if c.IsDebug {
		e.Use(middleware.Logger())
	}

	server := &Server{
		e:      e,
		config: c,
	}
	return server
}

func checkHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

func (s *Server) Start() error {
	s.e.Server.Addr = fmt.Sprintf(":%d", s.config.Server.Port)
	s.e.GET("/health", checkHealth)
	timeout := time.Second * 10
	return graceful.ListenAndServe(s.e.Server, timeout)
}
