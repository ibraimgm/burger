package telnet

import (
	"context"
	"fmt"
	"github.com/ibraimgm/burger/metrics"
	"github.com/ibraimgm/burger/recipe"
	"github.com/reiver/go-telnet"
	"github.com/reiver/go-telnet/telsh"
	"io"
	"net"

	"github.com/ibraimgm/burger/app"
)

type server struct {
	host    app.HostService
	collector *metrics.Collector

	srv      telnet.Server
	listener net.Listener
	doneCh   chan struct{}
	stopCh   chan struct{}
}

// New returns a new telnet server to the specified addr
func New(addr string, host app.HostService, collector *metrics.Collector) app.Server {
	return &server{
		srv:  telnet.Server{Addr: addr, Handler: nil},
		host: host,
		collector: collector,
	}
}

func (s *server) Start(ctx context.Context) {
	// invalid context
	select {
	case <-ctx.Done():
		return
	default:
	}

	// already running
	if s.listener != nil {
		return
	}

	// open telnet port
	listener, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		panic(err)
	}

	// server loop
	go func() {
		s.listener = listener
		s.stopCh = make(chan struct{})
		s.doneCh = make(chan struct{})

		// build shell commands
		shell := telsh.NewShellHandler()
		shell.WelcomeMessage = "Welcome to the BurgerNet(c) telnet server!\nType 'help' to see the list of commands.\n\n"

		shell.MustRegister("help", telsh.ProducerFunc(s.help))
		shell.MustRegister("status", telsh.ProducerFunc(s.status))

		shell.MustRegister("burger", telsh.ProducerFunc(s.burger))
		shell.MustRegister("b", telsh.ProducerFunc(s.burger))

		shell.MustRegister("doubleburger", telsh.ProducerFunc(s.doubleburger))
		shell.MustRegister("d", telsh.ProducerFunc(s.doubleburger))

		shell.MustRegister("hotdog", telsh.ProducerFunc(s.hotdog))
		shell.MustRegister("h", telsh.ProducerFunc(s.hotdog))

		s.srv.Handler = shell

		fmt.Printf("Starting telnet on address '%s'...\n", s.srv.Addr)
		if err := s.srv.Serve(listener); err != nil {
			fmt.Println(err)
		}

		close(s.stopCh)
		close(s.doneCh)
	}()

	// stop/cleanup routine
	go func() {
		select {
		case <-ctx.Done():
		case <-s.stopCh:
		}

		if err := s.listener.Close(); err != nil {
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

// telnet commands
func (s *server) help(c telnet.Context, namex string, argsx ...string) telsh.Handler {
	return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {
		fmt.Fprintf(stdout, "Available commands:\n")
		fmt.Fprintf(stdout, " - help         : Display this help message.\n")
		fmt.Fprintf(stdout, " - status       : Show restaurant status.\n")
		fmt.Fprintf(stdout, " - burger       : Ask for a new burger.\n")
		fmt.Fprintf(stdout, " - doubleburger : Ask for a Double burger.\n")
		fmt.Fprintf(stdout, " - hotdog       : Ask for a hotdog.\n\n")

		fmt.Fprintf(stdout, "You can use 'b', 'd' and 'h' as shortcuts for the 'burger', 'doubleburger' and 'hotdog', respectively.\n\n")

		return nil
	})
}

func (s *server) status(telnet.Context, string, ...string) telsh.Handler {
	return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {

		rc := s.collector.RecipeCount()
		served := s.collector.ServedByTable()
		line := s.collector.MaxLineSize()

		fmt.Fprintf(stdout, "=== Served Recipe Count ===\n")
		for k, v := range rc {
			fmt.Fprintf(stdout, "- %-15s ==> %d\n", k, v)
		}

		fmt.Fprintf(stdout, "\n")
		fmt.Fprintf(stdout, "=== Usage by Table ===\n")
		for k, v := range served {
			fmt.Fprintf(stdout, "- Table #%d ==> %d uses\n", k, v)
		}

		fmt.Fprintf(stdout, "\nMaximun line length: %d\n\n", line)

		return nil
	})
}

func (s *server) burger(telnet.Context, string, ...string) telsh.Handler {
	return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {

		s.host.Reserve(recipe.Burger)
		return nil
	})
}

func (s *server) doubleburger(telnet.Context, string, ...string) telsh.Handler {
	return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {

		s.host.Reserve(recipe.DoubleBurger)
		return nil
	})
}

func (s *server) hotdog(telnet.Context, string, ...string) telsh.Handler {
	return telsh.PromoteHandlerFunc(func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args ...string) error {

		s.host.Reserve(recipe.HotDog)
		return nil
	})
}
