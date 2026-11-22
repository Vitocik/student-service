// Package implementation реализует бизнес-логику сервиса студентов.
package implementation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"

	"example.com/student-service/student"
)

type service struct {
	repository student.Repository
	logger     log.Logger
}

// NewService создаёт сервис, работающий с переданным репозиторием.
func NewService(repo student.Repository, logger log.Logger) student.Service {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &service{repository: repo, logger: logger}
}

func (s *service) Create(ctx context.Context, st student.Student) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	st.FullName = strings.TrimSpace(st.FullName)
	st.Group = strings.TrimSpace(st.Group)
	if st.FullName == "" || st.Group == "" {
		return "", student.ErrInvalidStudent
	}

	st.ID = uuid.NewString()
	if st.EnrolledOn <= 0 {
		st.EnrolledOn = time.Now().Unix()
	}
	if err := s.repository.CreateStudent(ctx, st); err != nil {
		level.Error(logger).Log("err", err)
		return "", fmt.Errorf("%w: %v", student.ErrCmdRepository, err)
	}
	level.Info(logger).Log("student_id", st.ID, "msg", "student created")
	return st.ID, nil
}

func (s *service) GetByID(ctx context.Context, id string) (student.Student, error) {
	logger := log.With(s.logger, "method", "GetByID")
	id = strings.TrimSpace(id)
	if id == "" {
		return student.Student{}, student.ErrInvalidStudent
	}

	st, err := s.repository.GetStudentByID(ctx, id)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			return student.Student{}, student.ErrStudentNotFound
		}
		level.Error(logger).Log("err", err)
		return student.Student{}, fmt.Errorf("%w: %v", student.ErrQueryRepository, err)
	}
	return st, nil
}

func (s *service) ChangeGroup(ctx context.Context, id string, group string) error {
	logger := log.With(s.logger, "method", "ChangeGroup")
	id = strings.TrimSpace(id)
	group = strings.TrimSpace(group)
	if id == "" || group == "" {
		return student.ErrInvalidStudent
	}

	if err := s.repository.ChangeStudentGroup(ctx, id, group); err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			return student.ErrStudentNotFound
		}
		level.Error(logger).Log("err", err)
		return fmt.Errorf("%w: %v", student.ErrCmdRepository, err)
	}
	level.Info(logger).Log("student_id", id, "group", group, "msg", "student group changed")
	return nil
}
