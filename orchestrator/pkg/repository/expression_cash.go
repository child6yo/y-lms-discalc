package repository

import "github.com/child6yo/y-lms-discalc/orchestrator"

type expressionElement struct {
	expression string
	result     *orchestrator.Expression
}

func (c *expressionCache) Put(result *orchestrator.Expression) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[result.Expression]; ok {
		element.Value.(*expressionElement).result = result
		c.list.MoveToFront(element)
	} else {
		entry := &expressionElement{expression: result.Expression, result: result}
		element := c.list.PushFront(entry)
		c.cache[result.Expression] = element

		if c.list.Len() > c.capacity {
			oldest := c.list.Back()
			if oldest != nil {
				delete(c.cache, oldest.Value.(*expressionElement).expression)
				c.list.Remove(oldest)
			}
		}
	}
}

func (c *expressionCache) Get(expression string) (*orchestrator.Expression, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[expression]; ok {
		c.list.MoveToFront(element)
		return element.Value.(*expressionElement).result, true
	}

	return nil, false
}
