package service

import "gitlab.com/mfcekirdek/in-memory-store/internals/repository"

type StoreService interface {
}

type storeService struct {
	repository *repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	service := &storeService{repository: &repo}
	return service
}
