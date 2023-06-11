package app

import "github.com/inarithefox/partsy/server/config"

type AppOption func(a *App)
type AppOptionCreator func() []AppOption
type Option func(s *Server) error

func ServerConnector(p *Parts) AppOption {
	return func(a *App) {
		a.p = p
	}
}

func ConfigStore(configStore *config.Store) Option {
	return func(s *Server) error {
		s.ConfigStore(configStore)

		return nil
	}
}
