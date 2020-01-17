package http

import (
	"context"
	"fmt"
	"github.com/ibraimgm/burger/app"
	"github.com/ibraimgm/burger/metrics"
	"net"
	"net/http"
)

type server struct {
	collector *metrics.Collector
	srv       http.Server
	listener  net.Listener
	stopCh    chan struct{}
	doneCh    chan struct{}
}

func New(addr string, collector *metrics.Collector) app.Server {
	return &server{
		collector: collector,
		srv:       http.Server{Addr: addr},
	}
}

func (s *server) Start(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	if s.listener != nil {
		return
	}

	listener, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		panic(err)
	}

	go func() {
		s.listener = listener
		s.stopCh = make(chan struct{})
		s.doneCh = make(chan struct{})

		s.srv.Handler = s

		fmt.Printf("Starting http server on address '%s'...\n", s.srv.Addr)
		if err := s.srv.Serve(listener); err != nil {
			fmt.Println(err)
		}

		close(s.stopCh)
		close(s.doneCh)
	}()

	go func() {
		select {
		case <-ctx.Done():
		case <-s.stopCh:
		}

		if err := s.srv.Close(); err != nil {
			fmt.Println(err)
		}
	}()
}

func (s *server) Stop() {
	s.stopCh <- struct{}{}
}

func (s *server) Wait() {
	<-s.doneCh
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	rc := s.collector.RecipeCount()
	served := s.collector.ServedByTable()
	line := s.collector.MaxLineSize()

	fmt.Fprintf(w, "=== Served Recipe Count ===\n")
	for k, v := range rc {
		fmt.Fprintf(w, "- %-15s ==> %d\n", k, v)
	}

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "=== Usage by Table ===\n")
	for k, v := range served {
		fmt.Fprintf(w, "- Table #%d ==> %d uses\n", k, v)
	}

	fmt.Fprintf(w, "\nMaximun line length: %d\n\n", line)
}
