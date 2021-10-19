package repository

import "gitlab.com/mfcekirdek/in-memory-store/internals/model"

type StoreRepository interface {
}

type storeRepository struct {
	Store *model.Store
}

func NewStoreRepository(s *model.Store) StoreRepository {
	repository := &storeRepository{Store: s}
	return repository
}

func LoadStoreDataFromFile(filePath string) *model.Store {
	store := &model.Store{
		Pairs: []map[string]string{
			{"a": "b"},
		},
	}
	return store
}
