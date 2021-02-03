package servers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saaa-api/internal/core/config"

	"github.com/sirupsen/logrus"
)

// Server server
type Server struct {
	Server *http.Server
}

// NewServer create new server
func NewServer(port string, routes http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: routes,
		},
	}
}

// ListenAndServeWithGracefulShutdown listen and serve with graceful
func (s *Server) ListenAndServeWithGracefulShutdown() {
	trigger := make(chan os.Signal, 1)
	signal.Notify(trigger, os.Interrupt, syscall.SIGTERM)

	go s.ListenAndServe()

	// graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	cancel()

	srvCtx, srvCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer srvCancel()
	logrus.Infof("shutting down http server on port: %d ...", config.CF.App.Port)
	if err := s.Server.Shutdown(srvCtx); err != nil {
		logrus.Panicln("http server shutdown with error:", err)
	}
}

// ListenAndServe listen and serve
func (s *Server) ListenAndServe() {
	logrus.Infof("start http server on port: %d ...", config.CF.App.Port)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Panicf("listen: %s\n", err)
	}
}
