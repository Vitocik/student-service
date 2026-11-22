package implementation

import (
	"context"
	"testing"

	"github.com/go-kit/log"

	"example.com/student-service/repository"
	"example.com/student-service/student"
)

func TestCreateGetAndChangeGroup(t *testing.T) {
	svc := NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	ctx := context.Background()
	id, err := svc.Create(ctx, student.Student{FullName: "Иванов Иван", Group: "ПИ-23", GPA: 4.6})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if id == "" {
		t.Fatal("Create() returned empty ID")
	}

	created, err := svc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if created.FullName != "Иванов Иван" || created.Group != "ПИ-23" || created.GPA != 4.6 {
		t.Fatalf("unexpected student data: %+v", created)
	}
	if created.EnrolledOn <= 0 {
		t.Fatalf("EnrolledOn must be set, got %d", created.EnrolledOn)
	}

	if err := svc.ChangeGroup(ctx, id, "ПИ-24"); err != nil {
		t.Fatalf("ChangeGroup() error = %v", err)
	}
	changed, err := svc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() after ChangeGroup error = %v", err)
	}
	if changed.Group != "ПИ-24" {
		t.Fatalf("group was not changed: got %q", changed.Group)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	svc := NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	_, err := svc.GetByID(context.Background(), "unknown-id")
	if err != student.ErrStudentNotFound {
		t.Fatalf("GetByID() error = %v, want %v", err, student.ErrStudentNotFound)
	}
}

func TestCreateRejectsMissingNameOrGroup(t *testing.T) {
	svc := NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	_, err := svc.Create(context.Background(), student.Student{FullName: "", Group: "ПИ-23"})
	if err != student.ErrInvalidStudent {
		t.Fatalf("Create() error = %v, want %v", err, student.ErrInvalidStudent)
	}
}
