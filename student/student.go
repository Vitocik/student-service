// Package student содержит предметную модель и контракты сервиса студентов.
package student

import (
	"context"
	"errors"
)

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrInvalidStudent  = errors.New("invalid student data")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrCmdRepository   = errors.New("unable to update student repository")
	ErrQueryRepository = errors.New("unable to query student repository")
)

// Student — сущность студента из варианта 3.
// EnrolledOn хранит время зачисления в формате Unix timestamp (секунды).
type Student struct {
	ID         string  `json:"id,omitempty"`
	FullName   string  `json:"full_name"`
	Group      string  `json:"group"`
	GPA        float32 `json:"gpa"`
	EnrolledOn int64   `json:"enrolled_on,omitempty"`
}

// Service описывает бизнес-операции над студентами.
type Service interface {
	Create(ctx context.Context, student Student) (string, error)
	GetByID(ctx context.Context, id string) (Student, error)
	ChangeGroup(ctx context.Context, id string, group string) error
}

// Repository отделяет бизнес-логику от способа хранения данных.
type Repository interface {
	CreateStudent(ctx context.Context, student Student) error
	GetStudentByID(ctx context.Context, id string) (Student, error)
	ChangeStudentGroup(ctx context.Context, id string, group string) error
}
