package song

import "sync"

var nextId int64 = 0
var mux sync.Mutex

type Song struct {
	Id   int64
	Name string
}

func (s *Song) GetId() int64 {
	return s.Id
}

func (s *Song) GetName() string {
	return s.Name
}

func New(name string) *Song {
	s := new(Song)
	mux.Lock()
	s.Id = nextId
	nextId++
	mux.Unlock()
	s.Name = name

	return s
}