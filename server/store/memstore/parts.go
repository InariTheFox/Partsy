package memstore

import (
	"github.com/inarithefox/partsy/server/public/logger"
	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/store"
)

type MemPartStore struct {
	*MemStore

	parts []*model.Part
}

func NewMemPartStore(memStore *MemStore) store.PartStore {
	s := &MemPartStore{
		MemStore: memStore,
	}

	logger.Debug("new in-memory part store initialized")

	return s
}

func (s *MemPartStore) GetAllParts(page, pageSize int) (model.PartList, error) {
	r := s.parts

	return r, nil
}

func (s *MemPartStore) GetAllPartsCount() (int64, error) {
	count := (int64)(len(s.parts))

	return count, nil
}

func (s *MemPartStore) GetPart(partId string) (*model.Part, error) {
	p := &model.Part{
		Id: partId,
	}

	return p, nil
}

func (s *MemPartStore) Save(part *model.Part) (*model.Part, error) {
	if !model.IsValidId(part.Id) {
		part.Id = model.NewId()
		a := append(s.parts, part)
		s.parts = make([]*model.Part, len(a))
		copy(s.parts, a)
	} else {
		for i, v := range s.parts {
			if v.Id == part.Id {
				s.parts[i] = part
			}
		}
	}

	return part, nil
}
