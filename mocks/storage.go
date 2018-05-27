package mocks

// Storage is a mock
type Storage struct {
	Counter uint64
	Content map[string]string
}

// NewStorage returns a storage mock
func NewStorage() *Storage {
	return &Storage{0, map[string]string{}}
}

// SetPad update storage content in storage mock
func (s *Storage) SetPad(name string, value string) error {
	s.Content[name] = value
	return nil
}

// GetPad returns a content from storage mock
func (s *Storage) GetPad(name string) (value string, err error) {
	return s.Content[name], nil
}

// GetNextCounter returns next number of counter from storage mock
func (s *Storage) GetNextCounter() (next uint64, err error) {
	s.Counter = s.Counter + 1
	return s.Counter, nil
}
