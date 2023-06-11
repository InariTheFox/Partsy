package config

import (
	"encoding/json"
	"sync"

	"github.com/inarithefox/partsy/server/public/model"
	"github.com/inarithefox/partsy/server/public/utils"
	"github.com/pkg/errors"
)

var (
	ErrReadOnlyStore         = errors.New("configuration store is read-only")
	ErrReadOnlyConfiguration = errors.New("configuration is read-only")
)

type Store struct {
	backingStore BackingStore

	configLock sync.RWMutex
	config     *model.Config

	readOnly bool
}

type BackingStore interface {
	Set(*model.Config) error
	Load() ([]byte, error)
	GetFile(name string) ([]byte, error)
	SetFile(name string, data []byte) error
	HasFile(name string) (bool, error)
	RemoveFile(name string) error
	String() string
	Close() error
}

func NewStoreFromBacking(backingStore BackingStore, customDefaults *model.Config, readOnly bool) (*Store, error) {
	store := &Store{
		backingStore: backingStore,
		readOnly:     readOnly,
	}

	if err := store.Load(); err != nil {
		return nil, errors.Wrap(err, "unable to load on store creation")
	}

	return store, nil
}

func NewStoreFromDSN(dsn string, readOnly bool, customDefaults *model.Config, createFileIfNotExist bool) (*Store, error) {
	var err error
	var backingStore BackingStore

	if IsDatabaseDSN(dsn) {
		//backingStore, err = NewDatabaseStore(dsn)
	} else {
		backingStore, err = NewFileStore(dsn, createFileIfNotExist)
	}

	if err != nil {
		return nil, err
	}

	store, err := NewStoreFromBacking(backingStore, customDefaults, readOnly)
	if err != nil {
		backingStore.Close()
		return nil, errors.Wrap(err, "failed to create store")
	}

	return store, nil
}

func (s *Store) Close() error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	return s.backingStore.Close()
}

func (s *Store) Get() *model.Config {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return s.config
}

func (s *Store) GetFile(name string) ([]byte, error) {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return s.backingStore.GetFile(name)
}

func (s *Store) HasFile(name string) (bool, error) {
	s.configLock.RLock()
	defer s.configLock.RUnlock()

	return s.backingStore.HasFile(name)
}

func (s *Store) IsReadOnly() bool {
	s.configLock.RLock()
	defer s.configLock.RUnlock()

	return s.readOnly
}

func (s *Store) Load() error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	oldCfg := &model.Config{}
	if s.config != nil {
		oldCfg = s.config.Clone()
	}

	configBytes, err := s.backingStore.Load()
	if err != nil {
		return err
	}

	loadedCfg := &model.Config{}
	if len(configBytes) != 0 {
		if err = json.Unmarshal(configBytes, &loadedCfg); err != nil {
			return utils.HumanizeJSONError(err, configBytes)
		}
	}

	if loadedCfg.ServiceSettings.SiteURL == nil {
		loadedCfg.ServiceSettings.SiteURL = model.NewString("")
	}

	loadedCfg.SetDefaults()

	// TODO: Validate
	//if appErr := loadedCfg.IsValid(); appErr != nil {
	//	return errors.Wrap(appErr, "Invalid configuration.")
	//}

	hasChanged, err := equal(oldCfg, loadedCfg)
	if err != nil {
		return errors.Wrap(err, "failed to compare configurations")
	}

	if !s.readOnly && (hasChanged || len(configBytes) == 0) {
		err := s.backingStore.Set(loadedCfg)
		if err != nil && !errors.Is(err, ErrReadOnlyConfiguration) {
			return errors.Wrap(err, "failed to persist")
		}
	}

	s.config = loadedCfg

	if hasChanged {
		s.configLock.Unlock()
		s.configLock.Lock()
	}

	return nil
}

func (s *Store) RemoveFile(name string) error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	if s.readOnly {
		return ErrReadOnlyStore
	}

	return s.backingStore.RemoveFile(name)
}

func (s *Store) Set(newCfg *model.Config) (*model.Config, *model.Config, error) {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	if s.readOnly {
		return nil, nil, ErrReadOnlyStore
	}

	newCfg = newCfg.Clone()
	oldCfg := s.config.Clone()

	newCfg.SetDefaults()

	// TODO: Validate new configuration
	//if err := newCfg.IsValid(); err != nil {
	//	return nil, nil, errors.Wrap(err, "new configuration is invalid")
	//}

	if err := s.backingStore.Set(newCfg); err != nil {
		return nil, nil, errors.Wrap(err, "failed to persist store")
	}

	hasChanged, err := equal(oldCfg, newCfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to compare configurations")
	}

	s.config = newCfg
	newCfgCopy := newCfg.Clone()

	if hasChanged {
		s.configLock.Unlock()
		s.configLock.Lock()
	}

	return oldCfg, newCfgCopy, nil
}

func (s *Store) SetFile(name string, data []byte) error {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	if s.readOnly {
		return ErrReadOnlyStore
	}

	return s.backingStore.SetFile(name, data)
}

func (s *Store) String() string {
	return s.backingStore.String()
}
