package repository

import (
	"testing"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func TestCache1(t *testing.T) {
	cache := newExpressionCache(2)

	exp1 := &orchestrator.Expression{Result: 4, Expression: "2+2"}
	exp2 := &orchestrator.Expression{Result: 6, Expression: "3+3"}
	cache.Put(exp1)
	cache.Put(exp2)

	ch, ok := cache.Get("2+2")
	if !ok {
		t.Errorf("expected ok, got %v, false", ch)
	}
	if ch.Result != 4 {
		t.Errorf("expected 4, got %f", ch.Result)
	}

	exp3 := &orchestrator.Expression{Result: 8, Expression: "4+4"}
	cache.Put(exp3)

	ch, ok = cache.Get("3+3")
	if ok {
		t.Errorf("expected false, got %v, ok", ch)
	}
	if ch != nil {
		t.Errorf("expected nil, got %v", ch)
	}

	exp4 := &orchestrator.Expression{Result: 10, Expression: "5+5"}
	cache.Put(exp4)

	ch, ok = cache.Get("2+2")
	if ok {
		t.Errorf("expected false, got %v, ok", ch)
	}
	if ch != nil {
		t.Errorf("expected nil, got %v", ch)
	}

	ch, ok = cache.Get("4+4")
	if !ok {
		t.Errorf("expected ok, got %v, false", ch)
	}
	if ch.Result != 8 {
		t.Errorf("expected 4, got %f", ch.Result)
	}

	ch, ok = cache.Get("5+5")
	if !ok {
		t.Errorf("expected ok, got %v, false", ch)
	}
	if ch.Result != 10 {
		t.Errorf("expected 4, got %f", ch.Result)
	}
}
