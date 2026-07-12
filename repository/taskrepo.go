package repository

import (
	"errors"
	"server/domain"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrEmptyTitle = errors.New("Empty title")
var InternalServerError = errors.New("Enternal Server Error")

type InMemoryStorage struct {
	tasks map[int]*domain.Task
	mu    *sync.RWMutex
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[int]*domain.Task),
		mu:    &sync.RWMutex{},
	}
}

func (ms *InMemoryStorage) SaveTask(task *domain.Task) (*domain.Task, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.tasks[task.ID] = task
	return task, nil
}

func (ms *InMemoryStorage) GetAllTasks() []*domain.Task {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	arr := []*domain.Task{}
	for _, v := range ms.tasks {
		arr = append(arr, v)
	}
	return arr
}

func (ms *InMemoryStorage) GetTaskByID(id int) (*domain.Task, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if _, exists := ms.tasks[id]; exists {
		return ms.tasks[id], nil
	}
	return nil, ErrTaskNotFound

}
