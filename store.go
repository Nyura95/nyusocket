package nyusocket

import "errors"

// Store save the data client...
type Store struct {
	data map[string]interface{}
}

// NewStore create a new store
func NewStore() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

// SetData set a new value into the store
func (s *Store) SetData(key string, value interface{}) {
	s.data[key] = value
}

// GetData get data into the store
func (s *Store) GetData(key string) (interface{}, error) {
	if val, ok := s.data[key]; ok {
		return val, nil
	}
	return "", errors.New("key not exist")
}
