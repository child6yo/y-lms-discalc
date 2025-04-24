package repository

import (
	"container/list"
	"sync"
)

type Cache struct {
	cache     map[string]*list.Element
	capacity int
	list     *list.List
	mutex    sync.Mutex
}

func newCache(capacity int) *Cache {
	return &Cache{cache: map[string]*list.Element{}, capacity: capacity, list: list.New()}
}
