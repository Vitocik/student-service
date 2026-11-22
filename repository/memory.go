// Package repository содержит in-memory реализацию хранилища студентов.
package repository

import (
	"context"
	"strings"
	"sync"

	"example.com/student-service/student"
)

type memoryRepository struct {
	mu       sync.RWMutex
	students map[string]student.Student
}

// NewMemoryRepository создаёт потокобезопасное хранилище в памяти.
func NewMemoryRepository() student.Repository {
	return &memoryRepository{students: make(map[string]student.Student)}
}

func (r *memoryRepository) CreateStudent(ctx context.Context, st student.Student) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.students[st.ID] = st
	return nil
}

func (r *memoryRepository) GetStudentByID(ctx context.Context, id string) (student.Student, error) {
	if err := ctx.Err(); err != nil {
		return student.Student{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	st, ok := r.students[strings.TrimSpace(id)]
	if !ok {
		return student.Student{}, student.ErrStudentNotFound
	}
	return st, nil
}

func (r *memoryRepository) ChangeStudentGroup(ctx context.Context, id string, group string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	st, ok := r.students[strings.TrimSpace(id)]
	if !ok {
		return student.ErrStudentNotFound
	}
	st.Group = strings.TrimSpace(group)
	r.students[st.ID] = st
	return nil
}
