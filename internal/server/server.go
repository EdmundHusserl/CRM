package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/EdmundHusserl/CRM/internal/handlers"
	"github.com/EdmundHusserl/CRM/internal/repository"
	"github.com/EdmundHusserl/CRM/internal/repository/providers"
	"github.com/EdmundHusserl/CRM/internal/router"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Addr   string
	DB     repository.CustomerRepository
	Logger *logrus.Logger
	Router *mux.Router
}

func NewServer(repositoryProvider string, port int) Server {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	repo := providers.NewRepository(logger, repositoryProvider)
	handler := handlers.NewCustomerHandler(logger, repo)
	router := router.NewRouter(handler)

	return Server{
		Addr:   fmt.Sprintf(":%v", port),
		DB:     repo,
		Logger: logger,
		Router: router,
	}
}

func (s *Server) Listen() error {
	s.Logger.WithField(
		"event", fmt.Sprintf("Listening of port %v", s.Addr[1:]),
	).Info("Start server")
	return http.ListenAndServe(s.Addr, s.Router)
}
