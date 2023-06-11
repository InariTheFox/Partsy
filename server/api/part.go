package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/web"
)

func (api *Api) InitParts() {
	api.BaseRoutes.Parts.Handle("", api.APISessionRequired(getAllParts)).Methods("GET")
	api.BaseRoutes.Parts.Handle("", api.APISessionRequired(createPart)).Methods("POST")
}

func getAllParts(c *web.Context, w http.ResponseWriter, r *http.Request) {

	parts, err := c.App.GetAllParts(c.AppContext, c.Params.Page, c.Params.PageSize)
	if err != nil {
		c.Err = err
		return
	}

	if c.Params.IncludeTotalCount {
		totalCount, err := c.App.GetAllPartsCount(c.AppContext)
		if err != nil {
			c.Err = err
			return
		}

		pwc := &model.PartsWithCount{
			Parts:      parts,
			TotalCount: totalCount,
		}

		if err := json.NewEncoder(w).Encode(pwc); err != nil {
			logger.Warn(fmt.Sprintf("error while writing response: %s", err))
		}
	}

	if err := json.NewEncoder(w).Encode(parts); err != nil {
		logger.Warn(fmt.Sprintf("error while writing response: %s", err))
	}
}

func createPart(c *web.Context, w http.ResponseWriter, r *http.Request) {

}
