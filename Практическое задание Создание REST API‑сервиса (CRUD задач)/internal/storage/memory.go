package storage

import (
	"errors"
	"sync"
	"time"

	"tasks-api/internal/models"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (m *MemoryStorage) List() []models.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]models.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result
}

func (m *MemoryStorage) Create(task models.Task) (models.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if task.Title == "" {
		return models.Task{}, errors.New("title is required")
	}

	task.ID = m.nextID
	m.nextID++
	task.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	m.tasks[task.ID] = task
	return task, nil
}

func (m *MemoryStorage) Get(id int) (models.Task, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[id]
	return task, exists
}

func (m *MemoryStorage) Update(id int, task models.Task) (models.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[id]; !exists {
		return models.Task{}, errors.New("task not found")
	}

	if task.Title == "" {
		return models.Task{}, errors.New("title is required")
	}

	task.ID = id
	if existing, ok := m.tasks[id]; ok {
		task.CreatedAt = existing.CreatedAt
	}

	m.tasks[id] = task
	return task, nil
}

func (m *MemoryStorage) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(m.tasks, id)
	return nil
}
