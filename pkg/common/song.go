package common

import "sync"

var nextId int64 = 0
var mux sync.Mutex

type Song struct {
	id int64
	name string
}

func New(name string) *Song {
	s := new(Song)
	mux.Lock()
	s.id = nextId
	nextId++
	mux.Unlock()
	s.name = name

	return s
}