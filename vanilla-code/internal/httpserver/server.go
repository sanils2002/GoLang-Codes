package httpserver

import (
	"log"

	"vanilla-code/internal/auth"
	"vanilla-code/internal/config"
	"vanilla-code/internal/user"
)

type Server struct {
	cfg      config.HTTPConfig
	authSvc  *auth.Service
	userSvc  *user.Service
}

func New(cfg config.HTTPConfig, authSvc *auth.Service, userSvc *user.Service) *Server {
	return &Server{cfg: cfg, authSvc: authSvc, userSvc: userSvc}
}

func (s *Server) Run() {
	log.Printf("http server (stub) would listen on %s:%d", s.cfg.Host, s.cfg.Port)
	_ = s.authSvc
	_ = s.userSvc
}
