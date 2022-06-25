package nyusocket

import "errors"

// Store save the data client...
type Store struct {
	data  map[string]interface{}
	sdata map[string]string
}

// NewStore create a new store
func NewStore() Store {
	return Store{
		data:  make(map[string]interface{}),
		sdata: make(map[string]string),
	}
}

// SetData set a new value into the store
func (s *Store) SetData(key string, value interface{}) {
	s.data[key] = value
}

// GetData get data into the store
func (s Store) GetData(key string) (interface{}, error) {
	if val, ok := s.data[key]; ok {
		return val, nil
	}
	return "", errors.New("key not exist")
}

func (s *Store) SetSData(key, value string) {
	s.sdata[key] = value
}

func (s Store) GetSData(key string) string {
	if val, ok := s.sdata[key]; ok {
		return val
	}
	return ""
}
