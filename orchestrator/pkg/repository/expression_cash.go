package repository

import "github.com/child6yo/y-lms-discalc/orchestrator"

type ExpressionCache struct {
	expression string
	result     *orchestrator.Result
}

func (r *Repository) CacheResult(result *orchestrator.Result) {
	r.Cache.mutex.Lock()
	defer r.Cache.mutex.Unlock()

	if element, ok := r.Cache.cache[result.Expression]; ok {
		element.Value.(*ExpressionCache).result = result
		r.Cache.list.MoveToFront(element)
	} else {
		entry := &ExpressionCache{expression: result.Expression, result: result}
		element := r.Cache.list.PushFront(entry)
		r.Cache.cache[result.Expression] = element

		if r.Cache.list.Len() > r.Cache.capacity {
			oldest := r.Cache.list.Back()
			if oldest != nil {
				delete(r.Cache.cache, oldest.Value.(*ExpressionCache).expression)
				r.Cache.list.Remove(oldest)
			}
		}
	}
}

func (r *Repository) GetCachedResult(expression string) (*orchestrator.Result, bool) {
	r.Cache.mutex.Lock()
	defer r.Cache.mutex.Unlock()

	if element, ok := r.Cache.cache[expression]; ok {
		r.Cache.list.MoveToFront(element)
		return element.Value.(*ExpressionCache).result, true
	}

	return nil, false
}
