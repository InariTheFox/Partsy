package store

import (
	"context"

	"github.com/inarithefox/partsy/server/public/model"
)

type Store interface {
	Close()
	Context() context.Context
	Part() PartStore
	SetContext(context context.Context)
}

type PartStore interface {
	GetAllParts(page, pageSize int) (model.PartList, error)
	GetAllPartsCount() (int64, error)
	Save(part *model.Part) (*model.Part, error)
}
