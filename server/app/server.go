package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/inarithefox/partsy/server/config"
	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/public/utils"
	"github.com/inarithefox/partsy/server/store"
	"github.com/inarithefox/partsy/server/store/memstore"
	"github.com/pkg/errors"
)

type Server struct {
	RootRouter *mux.Router
	Router     *mux.Router

	Server     *http.Server
	ListenAddr *net.TCPAddr

	configStore *config.Store

	config          *model.Config
	didFinishListen chan struct{}
	p               *Parts
	store           store.Store
}

func NewServer(options ...Option) (*Server, error) {
	rootRouter := mux.NewRouter()

	s := &Server{
		RootRouter: rootRouter,
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	subpath, err := utils.GetSubpathFromConfig(s.Config())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse SiteURL subpath from configuration")
	}

	s.Router = s.RootRouter.PathPrefix(subpath).Subrouter()

	store := memstore.New()

	p, err := NewParts(store)
	if err != nil {
		logger.Error(err, "unable to initialize part infrastructure")
		return nil, errors.Wrap(err, "failed to initialize part infrastructure")
	}

	s.p = p
	s.store = store

	logger.Info(fmt.Sprintf("Fox Labs Partsy - Version: %v", model.CurrentVersion))
	logger.Info(fmt.Sprintf("GO Version: %v", runtime.Version()))

	return s, nil
}

func (s *Server) ConfigStore(configStore *config.Store) {
	logger.Debug("setting configuration store")
	s.configStore = configStore
	s.config = configStore.Get()
}

func (s *Server) Parts() *Parts {
	return s.p
}

func handleHTTPRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}

	target := "https://" + stripPort(r.Host) + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusFound)
}

func (s *Server) Shutdown() error {
	logger.Info("Stopping server...")
	if s.store != nil {
		s.Store().Close()
	}

	s.StopHTTPServer()

	return nil
}

func (s *Server) Start() error {

	logger.Info("Starting server...")

	var handler http.Handler = s.RootRouter

	s.Server = &http.Server{
		Handler:      handler,
		ReadTimeout:  time.Duration(*s.Config().ServiceSettings.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(*s.Config().ServiceSettings.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(*s.Config().ServiceSettings.IdleTimeout) * time.Second,
	}

	addr := *s.Config().ServiceSettings.ListenAddress
	if addr == "" {
		if *s.Config().ServiceSettings.ConnectionSecurity == model.ConnectionSecurityTLS {
			addr = ":https"
		} else {
			addr = ":http"
		}
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "critical error occured during start of server")
	}

	s.ListenAddr = listener.Addr().(*net.TCPAddr)

	logger.Info(fmt.Sprintf("Server is listening on %v", listener.Addr().String()))

	s.didFinishListen = make(chan struct{})
	go func() {
		var err error

		if *s.Config().ServiceSettings.ConnectionSecurity == model.ConnectionSecurityTLS {

		} else {
			err = s.Server.Serve(listener)
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err, "Error starting server")
			time.Sleep(time.Second)
		}

		close(s.didFinishListen)
	}()

	return nil
}

func stripPort(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	return net.JoinHostPort(host, "443")
}

func (s *Server) StopHTTPServer() {
	if s.Server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		didShutdown := false

		for s.didFinishListen != nil && !didShutdown {
			if err := s.Server.Shutdown(ctx); err != nil {
				logger.Error(err, "unable to shutdown server")
			}

			timer := time.NewTimer(time.Millisecond * 60)
			select {
			case <-s.didFinishListen:
				didShutdown = true
			case <-timer.C:
			}
		}

		s.Server.Close()
		s.Server = nil
	}
}

func (s *Server) Store() store.Store {
	return s.store
}
