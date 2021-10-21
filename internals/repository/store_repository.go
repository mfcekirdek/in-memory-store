package repository

type StoreRepository interface {
	Get(key string) string
	Set(key string, value string) bool
	Flush() map[string]string
	GetStore() map[string]string
	LoadStore(store map[string]string)
}

type storeRepository struct {
	store map[string]string
}

func NewStoreRepository() StoreRepository {
	repository := &storeRepository{}
	return repository
}

func (s storeRepository) GetStore() map[string]string {
	return s.store
}

func (s storeRepository) LoadStore(store map[string]string) {
	s.store = store
}

func (s storeRepository) Flush() map[string]string {
	s.store = map[string]string{}
	return s.store
}

func (s storeRepository) Get(key string) string {
	if value, ok := s.store[key]; ok {
		return value
	}
	return ""
}

func (s storeRepository) Set(key, value string) bool {
	keyAlreadyExists := false
	if s.store[key] != "" {
		keyAlreadyExists = true
	}
	s.store[key] = value
	return keyAlreadyExists
}
