package app

import "github.com/inarithefox/partsy/server/public/model"

func (a *App) Config() *model.Config {
	return a.p.srv.Config()
}

func (s *Server) Config() *model.Config {
	return s.configStore.Get()
}
