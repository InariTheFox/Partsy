package app

import "github.com/inarithefox/partsy/server/public/logger"

type App struct {
	srv *Server
}

func New(options ...AppOption) *App {
	app := &App{}

	for _, option := range options {
		option(app)
	}

	logger.Debug("new application initialized")

	return app
}

func (a *App) Srv() *Server {
	return a.srv
}
