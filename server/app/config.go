package app

import "github.com/inarithefox/partsy/server/public/model"

func (a *App) Config() *model.Config {
	return a.srv.Config()
}

func (s *Server) Config() *model.Config {
	return s.configStore.Get()
}
