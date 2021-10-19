package service

import (
	"gitlab.com/mfcekirdek/in-memory-store/internals/repository"
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
)

type StoreService interface {
	Get(key string) (map[string]string, error)
	Set(key string, value string) map[string]string
}

type storeService struct {
	repository repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	service := &storeService{repository: repo}
	return service
}

func (s *storeService) Get(key string) (map[string]string, error) {
	value := s.repository.Get(key)
	if value == "" {
		return nil, utils.ErrNotFound
	}
	return map[string]string{key: value}, nil
}

func (s *storeService) Set(key string, value string) map[string]string {
	s.repository.Set(key, value)
	return map[string]string{key: value}
}
