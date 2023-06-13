package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/web"
)

type GraphQLInput struct {
	Query         string         `json:"query"`
	OperationName string         `json:"oeprationName"`
	Variables     map[string]any `json:"variables"`
}

func (api *Api) InitGraphQL() error {

	api.BaseRoutes.ApiRoot.Handle("/graphql", api.APIHandler(graphiQL)).Methods("GET")
	api.BaseRoutes.ApiRoot.Handle("/graphql", api.APISessionRequired(api.graphQL)).Methods("POST")

	return nil
}

func (api *Api) graphQL(c *web.Context, w http.ResponseWriter, r *http.Request) {
	var response *graphql.Response
	defer func() {
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				logger.Warn(fmt.Sprintf("error while writing response: %v", err))
			}
		}
	}()

	r.Body = http.MaxBytesReader(w, r.Body, 102400)

	var params GraphQLInput
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Error(err, "invalid request body")
		innerErr := gqlerrors.Errorf("invalid request body: %v", err)
		response = &graphql.Response{Errors: []*gqlerrors.QueryError{innerErr}}
		return
	}

	c.GraphQLOperationName = params.OperationName

	rContext := r.Context()
	rContext = context.WithValue(rContext, 0, c)

	response = api.schema.Exec(
		rContext,
		params.Query,
		params.OperationName,
		params.Variables)

	if len(response.Errors) > 0 {
		logger.Error(nil, fmt.Sprintf("error executing operation: %s: %v", params.OperationName, response.Errors))
	}
}

func graphiQL(c *web.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	var page = []byte(`<!DOCTYPE html><html><head><title>Partsy GraphQL Editor</title></head><body>Editor coming soon!</body></html>`)

	w.Write(page)
}
