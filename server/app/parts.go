package app

import (
	"errors"
	"net/http"

	"github.com/inarithefox/partsy/server/app/request"
	"github.com/inarithefox/partsy/server/public/model"
)

type Parts struct {
	srv *Server
}

func NewParts(s *Server) (*Parts, error) {
	if s == nil {
		return nil, errors.New("server not passed to part infrastructure")
	}

	p := &Parts{
		srv: s,
	}

	s.p = p

	return p, nil
}

func (a *App) GetAllParts(c *request.Context, page, pageSize int) (model.PartList, *model.AppError) {
	parts, err := a.Srv().Store().Part().GetAllParts(page*pageSize, pageSize)
	if err != nil {
		return nil, model.NewAppError(("GetAllParts"), "app.parts.GetAllParts", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return parts, nil
}

func (a *App) GetAllPartsCount(c *request.Context) (int64, *model.AppError) {
	count, err := a.Srv().Store().Part().GetAllPartsCount()
	if err != nil {
		return 0, model.NewAppError("GetAllPartsCount", "app.parts.GetAllPartsCount", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return count, nil
}
