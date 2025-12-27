package config

import (
	"encoding/json"
)

type Storage struct {
	path           string // ro
	env            string // ro internal value
	saltSize       int    // ro internal value
	argonTime      uint32 // ro internal value
	argonMemory    uint32 // ro internal value
	argonKeyLength uint32 // ro internal value
	argonThreads   uint8  // ro internal value
}

type storage struct {
	Path string `json:"path"`
}

func newStorage() *Storage {
	return &Storage{
		path:           "./storage", // default
		env:            "WBLK",
		saltSize:       32,
		argonTime:      3,
		argonMemory:    64 * 1024,
		argonThreads:   4,
		argonKeyLength: 32,
	}
}

func (s *Storage) Path() string           { return s.path }
func (s *Storage) Env() string            { return s.env }
func (s *Storage) SaltSize() int          { return s.saltSize }
func (s *Storage) ArgonTime() uint32      { return s.argonTime }
func (s *Storage) ArgonMemory() uint32    { return s.argonMemory }
func (s *Storage) ArgonKeyLength() uint32 { return s.argonKeyLength }
func (s *Storage) ArgonThreads() uint8    { return s.argonThreads }

func (s *Storage) UnmarshalJSON(b []byte) error {
	temp := &storage{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	s.path = temp.Path
	s.env = "WBLK"
	s.saltSize = 32
	s.argonTime = 3
	s.argonMemory = 64 * 1024
	s.argonThreads = 4
	s.argonKeyLength = 32
	return nil
}

func (s *Storage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&storage{
		Path: s.path,
	})
}
