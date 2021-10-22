// Package repository executes the repository operations.
// This layer is used by service layer.
package repository

// StoreRepository interface has Get, Set, Flush, GetStore and LoadStore functions.
type StoreRepository interface {
	Get(key string) string
	Set(key string, value string) bool
	Flush() map[string]string
	GetStore() map[string]string
	LoadStore(store map[string]string)
}

// storeRepository implements StoreRepository Interface.
type storeRepository struct {
	store map[string]string
}

// NewStoreRepository creates a new storeRepository instance.
func NewStoreRepository() StoreRepository {
	repository := &storeRepository{}
	return repository
}

// Returns store contents
func (s *storeRepository) GetStore() map[string]string {
	return s.store
}

// Takes the <map> parameter and saves it to the store
func (s *storeRepository) LoadStore(store map[string]string) {
	s.store = store
}

// Deletes all items and returns empty store
func (s *storeRepository) Flush() map[string]string {
	s.store = map[string]string{}
	return s.store
}

// Takes <string> key parameter.
// If the key is found, returns it's value.
// Otherwise returns empty string.
func (s *storeRepository) Get(key string) string {
	if value, ok := s.store[key]; ok {
		return value
	}
	return ""
}

// Takes <string>,<string> key, value parameters.
// Sets value of key as given value.
// If the key already exists returns true.
// Otherwise returns false.
func (s *storeRepository) Set(key, value string) bool {
	keyAlreadyExists := false
	if s.store == nil {
		s.store = map[string]string{}
	}
	if s.store[key] != "" {
		keyAlreadyExists = true
	}
	s.store[key] = value
	return keyAlreadyExists
}
