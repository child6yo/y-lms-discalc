package repository

import "github.com/child6yo/y-lms-discalc/orchestrator"

// expressionElement - списочный элемент кэша.
type expressionElement struct {
	expression string
	result     *orchestrator.Expression
}

// Put добавляет арифметическое выражение в кэш.
// На вход принимает структуру арифметического выражения.
func (c *expressionCache) Put(expression *orchestrator.Expression) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[expression.Expression]; ok {
		element.Value.(*expressionElement).result = expression
		c.list.MoveToFront(element)
	} else {
		entry := &expressionElement{expression: expression.Expression, result: expression}
		element := c.list.PushFront(entry)
		c.cache[expression.Expression] = element

		if c.list.Len() > c.capacity {
			oldest := c.list.Back()
			if oldest != nil {
				delete(c.cache, oldest.Value.(*expressionElement).expression)
				c.list.Remove(oldest)
			}
		}
	}
}

// Get получает значение из кэша.
// На вход принимает ключ в виде арифметического выражения в строковом виде.
func (c *expressionCache) Get(expression string) (*orchestrator.Expression, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, ok := c.cache[expression]; ok {
		c.list.MoveToFront(element)
		return element.Value.(*expressionElement).result, true
	}

	return nil, false
}
