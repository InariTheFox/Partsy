package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-go"
	"github.com/inarithefox/partsy/server/app"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/web"
)

type Routes struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // api

	Parts          *mux.Router // api/parts
	Part           *mux.Router // api/parts/{part_id:[A-Za-z0-9]}
	PartCategories *mux.Router // api/parts/categories

	Files *mux.Router // api/files
	File  *mux.Router // api/files/{file_id:[A-Za-z0-9]}
}

type Api struct {
	srv        *app.Server
	schema     *graphql.Schema
	BaseRoutes *Routes
}

func Init(srv *app.Server) (*Api, error) {
	api := &Api{
		srv:        srv,
		BaseRoutes: &Routes{},
	}

	api.BaseRoutes.Root = srv.Router
	api.BaseRoutes.ApiRoot = srv.Router.PathPrefix(model.ApiURLSuffix).Subrouter()

	api.BaseRoutes.Parts = api.BaseRoutes.ApiRoot.PathPrefix("/parts").Subrouter()
	api.BaseRoutes.Part = api.BaseRoutes.Parts.PathPrefix("/{part_id:[A-Za-z0-9]+}").Subrouter()

	api.InitParts()

	srv.Router.Handle("/api/{anything:.*}", http.HandlerFunc(api.HandleNotFound))

	return api, nil
}

func (api *Api) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	app := app.New(app.ServerConnector(api.srv))
	web.HandleNotFound(app, w, r)
}
