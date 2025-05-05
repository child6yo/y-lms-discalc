package repository

import (
	"container/list"
	"sync"
)

// expressionCache реализует LRU кэш для хранения результатов вычисления арифметических выражений.
type expressionCache struct {
	cache    map[string]*list.Element
	capacity int
	list     *list.List
	mutex    sync.Mutex
}

// newExpressionCache создает новый LRU кэш для хранения результатов вычисления арифметических выражений.
//
// Параметры:
//   - capacity: вместимость LRU кеша.
func newExpressionCache(capacity int) *expressionCache {
	return &expressionCache{cache: map[string]*list.Element{}, capacity: capacity, list: list.New()}
}
