package app

type App struct {
	p *Parts
}

func New(options ...AppOption) *App {
	app := &App{}

	for _, option := range options {
		option(app)
	}

	return app
}

func (a *App) Srv() *Server {
	return a.p.srv
}
