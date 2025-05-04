package repository

import (
	"container/list"
	"sync"
)

type expressionCache struct {
	cache     map[string]*list.Element
	capacity int
	list     *list.List
	mutex    sync.Mutex
}

func newExpressionCache(capacity int) *expressionCache {
	return &expressionCache{cache: map[string]*list.Element{}, capacity: capacity, list: list.New()}
}
