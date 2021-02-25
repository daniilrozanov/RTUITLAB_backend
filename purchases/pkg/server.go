package templates

import (
	"net/http"
	"time"
)

type Server struct{
	http_serv *http.Server
}

func (s *Server) Start(port string, handler http.Handler) error {
	s.http_serv = &http.Server{
		Addr: ":" + port,
		Handler: handler,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
	}
	return (s.http_serv.ListenAndServe())
}