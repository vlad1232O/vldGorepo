package models

import (
	"sync"
	"time"
)

type Expense struct {
	Amount      float64
	Description string
	Date        time.Time
}

type ExpenseStore struct {
	mu       sync.RWMutex
	expenses []Expense
}

func NewExpenseStore() *ExpenseStore {
	return &ExpenseStore{}
}

func (s *ExpenseStore) Add(e Expense) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expenses = append(s.expenses, e)
}

func (s *ExpenseStore) All() []Expense {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Expense, len(s.expenses))
	copy(result, s.expenses)
	return result
}
