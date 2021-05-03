package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ServerInterface interface {
	Start() error
	Stop() error
}

type Server struct{
	*http.Server
	lis net.Listener
	network string
	address string
	timeout time.Duration
	router *mux.Router
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.timeout)
	defer cancel()
	s.router.ServeHTTP(res, req.WithContext(ctx))
}


func NewServer()*Server{
	srv := &Server{
		network: "tcp",
		address: "127.0.0.1:8087",
		timeout: time.Second,
	}
	srv.Server = &http.Server{Handler: srv}
	return srv
}
func (s *Server)Start()error{
	lis,err := net.Listen(s.network,s.address)
	if err != nil{
		return err
	}
	s.lis = lis
	if err := s.Serve(lis); !errors.Is(err, http.ErrServerClosed){
		return err
	}
	return nil
}
func (s *Server)Stop()error{
	return s.Shutdown(context.Background())
}

