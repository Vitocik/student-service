package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"example.com/student-service/student"
)

func TestMemoryRepositoryCreateGetAndChangeGroup(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()
	st := student.Student{ID: "student-1", FullName: "Иванов Иван", Group: "ПИ-23", GPA: 4.5, EnrolledOn: 1}
	if err := repo.CreateStudent(ctx, st); err != nil {
		t.Fatalf("CreateStudent() error = %v", err)
	}
	got, err := repo.GetStudentByID(ctx, st.ID)
	if err != nil {
		t.Fatalf("GetStudentByID() error = %v", err)
	}
	if got != st {
		t.Fatalf("GetStudentByID() = %+v, want %+v", got, st)
	}
	if err := repo.ChangeStudentGroup(ctx, st.ID, "ПИ-24"); err != nil {
		t.Fatalf("ChangeStudentGroup() error = %v", err)
	}
	got, err = repo.GetStudentByID(ctx, st.ID)
	if err != nil {
		t.Fatalf("GetStudentByID() after update error = %v", err)
	}
	if got.Group != "ПИ-24" {
		t.Fatalf("group = %q, want ПИ-24", got.Group)
	}
}

func TestMemoryRepositoryNotFound(t *testing.T) {
	_, err := NewMemoryRepository().GetStudentByID(context.Background(), "missing")
	if !errors.Is(err, student.ErrStudentNotFound) {
		t.Fatalf("error = %v, want ErrStudentNotFound", err)
	}
}

func TestMemoryRepositoryConcurrentAccess(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()
	const count = 80
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("student-%d", i)
			st := student.Student{ID: id, FullName: "Student", Group: "G-1"}
			if err := repo.CreateStudent(ctx, st); err != nil {
				t.Errorf("CreateStudent(%s): %v", id, err)
				return
			}
			if err := repo.ChangeStudentGroup(ctx, id, "G-2"); err != nil {
				t.Errorf("ChangeStudentGroup(%s): %v", id, err)
				return
			}
			got, err := repo.GetStudentByID(ctx, id)
			if err != nil {
				t.Errorf("GetStudentByID(%s): %v", id, err)
				return
			}
			if got.Group != "G-2" {
				t.Errorf("student %s group = %s, want G-2", id, got.Group)
			}
		}(i)
	}
	wg.Wait()
}

func TestMemoryRepositoryRespectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := NewMemoryRepository().CreateStudent(ctx, student.Student{ID: "student-1"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("CreateStudent() error = %v, want context.Canceled", err)
	}
}
