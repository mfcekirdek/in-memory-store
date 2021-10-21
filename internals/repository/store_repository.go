package repository

type IStoreRepository interface {
	Get(key string) string
	Set(key string, value string) bool
	Flush() map[string]string
	GetStore() map[string]string
	LoadStore(store map[string]string)
}

type StoreRepository struct {
	Store map[string]string
}

func NewStoreRepository() IStoreRepository {
	repository := &StoreRepository{}
	return repository
}

func (s StoreRepository) GetStore() map[string]string {
	return s.Store
}

func (s StoreRepository) LoadStore(store map[string]string) {
	s.Store = store
}

func (s StoreRepository) Flush() map[string]string {
	s.Store = map[string]string{}
	return s.Store
}

func (s StoreRepository) Get(key string) string {
	if value, ok := s.Store[key]; ok {
		return value
	}
	return ""
}

func (s StoreRepository) Set(key, value string) bool {
	keyAlreadyExists := false
	if s.Store[key] != "" {
		keyAlreadyExists = true
	}
	s.Store[key] = value
	return keyAlreadyExists
}
