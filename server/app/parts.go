package app

import (
	"errors"
	"net/http"

	"github.com/inarithefox/partsy/server/app/request"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/store"
)

type Parts struct {
	store store.PartStore
}

func NewParts(s store.Store) (*Parts, error) {
	if s == nil {
		return nil, errors.New("server not passed to part infrastructure")
	}

	p := &Parts{
		store: s.Part(),
	}

	return p, nil
}

func (a *App) CreatePart(c *request.Context, part *model.Part) (*model.Part, *model.AppError) {
	rp, err := a.Srv().Parts().Store().Save(part)
	if err != nil {
		return nil, model.NewAppError("CreatePart", "app.parts.create_part", nil, "", http.StatusBadRequest).Wrap(err)
	}

	return rp, nil
}

func (a *App) GetAllParts(c *request.Context, page, pageSize int) (model.PartList, *model.AppError) {
	parts, err := a.Srv().Parts().Store().GetAllParts(page*pageSize, pageSize)
	if err != nil {
		return nil, model.NewAppError("GetAllParts", "app.parts.get_all_parts", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return parts, nil
}

func (a *App) GetAllPartsCount(c *request.Context) (int64, *model.AppError) {
	count, err := a.Srv().Parts().Store().GetAllPartsCount()
	if err != nil {
		return 0, model.NewAppError("GetAllPartsCount", "app.parts.get_all_parts_count", nil, "", http.StatusInternalServerError).Wrap(err)
	}

	return count, nil
}

func (a *App) GetPart(c *request.Context, partId string) (*model.Part, *model.AppError) {
	part, err := a.Srv().Parts().Store().GetPart(partId)
	if err != nil {
		return nil, model.NewAppError("GetPart", "app.parts.get_part", nil, "", http.StatusNotFound)
	}

	return part, nil
}

func (p *Parts) Store() store.PartStore {
	return p.store
}
