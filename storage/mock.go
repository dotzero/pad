package storage

// MockStorage is a storage mock
type MockStorage struct {
	Counter uint64
	Content map[string]string
}

// NewMock returns a storage mock
func NewMock() *MockStorage {
	return &MockStorage{0, map[string]string{}}
}

// SetPad update storage content in storage mock
func (s *MockStorage) Set(name string, value string) error {
	s.Content[name] = value
	return nil
}

// GetPad returns a content from storage mock
func (s *MockStorage) Get(name string) (value string, err error) {
	return s.Content[name], nil
}

// GetNextCounter returns next number of counter from storage mock
func (s *MockStorage) NextCounter() (next uint64, err error) {
	s.Counter = s.Counter + 1
	return s.Counter, nil
}
