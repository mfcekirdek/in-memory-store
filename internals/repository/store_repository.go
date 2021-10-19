package repository

type StoreRepository interface {
}

type storeRepository struct {
}

func NewStoreRepository() StoreRepository {
	repository := &storeRepository{}
	return repository
}
