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

	api.BaseRoutes.Part.Handle("", api.APISessionRequired(getPart)).Methods("GET")
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
	} else {
		if err := json.NewEncoder(w).Encode(parts); err != nil {
			logger.Warn(fmt.Sprintf("error while writing response: %s", err))
		}
	}
}

func getPart(c *web.Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePartId()
	if c.Err != nil {
		return
	}

	part, err := c.App.GetPart(c.AppContext, c.Params.PartId)
	if err != nil {
		c.Err = err
		return
	}

	if err := json.NewEncoder(w).Encode(part); err != nil {
		logger.Warn(fmt.Sprintf("error while writing response: %v", err))
	}
}

func createPart(c *web.Context, w http.ResponseWriter, r *http.Request) {
	var part *model.Part
	err := json.NewDecoder(r.Body).Decode(&part)
	if err != nil {
		c.SetInvalidParameterWithErr("part", err)
		return
	}

	// TODO: Audit

	cp, appErr := c.App.CreatePart(c.AppContext, part)
	if appErr != nil {
		c.Err = appErr
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(cp); err != nil {
		logger.Warn(fmt.Sprintf("error while writing response: %v", err))
	}
}
