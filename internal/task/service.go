package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	Create(description string, dueDate time.Time) (*Task, error)
	GetByID(id uuid.UUID) (*Task, error)
	GetTaskByPartialId(id string) (*Task, error)
	List(selector *TaskSelector, filter *TaskFilter) ([]Task, error)
	Complete(id string) error
	Delete(id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) TaskService {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(description string, dueDate time.Time) (*Task, error) {
	task := &Task{
		Description: description,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}

	if err := task.Validate(); err != nil {
		return nil, &Error{Op: "Create", Err: err}
	}

	if err := s.repository.Save(task); err != nil {
		return nil, &Error{Op: "Create", Err: err}
	}

	return task, nil
}

func (s *service) GetByID(id uuid.UUID) (*Task, error) {
	task, err := s.repository.GetByID(id)
	if err != nil {
		return nil, &Error{Op: "GetByID", Err: err}
	}
	return task, nil
}

func (s *service) GetTaskByPartialId(id string) (*Task, error) {
	task, err := s.repository.GetTaskByPartialId(id)
	if err != nil {
		return nil, &Error{Op: "GetByID", Err: err}
	}
	return task, nil
}

func (s *service) List(selector *TaskSelector, filter *TaskFilter) ([]Task, error) {
	if selector == nil {
		selector = NewTaskSelector()
	}
	if filter == nil {
		filter = NewTaskFilter()
	}

	tasks, err := s.repository.List(selector, filter)
	if err != nil {
		return nil, &Error{Op: "List", Err: err}
	}
	return tasks, nil
}

func (s *service) Complete(id string) error {
	task, err := s.repository.GetTaskByPartialId(id)
	if err != nil {
		return &Error{Op: "Complete", Err: err}
	}

	task.IsCompleted = true
	if err := s.repository.Update(task); err != nil {
		return &Error{Op: "Complete", Err: err}
	}

	return nil
}

func (s *service) Delete(id string) error {
	task, err := s.repository.GetTaskByPartialId(id)
	if err != nil {
		return &Error{Op: "Delete", Err: err}
	}

	if err := s.repository.Delete(task); err != nil {
		return &Error{Op: "Delete", Err: err}
	}

	return nil
}
