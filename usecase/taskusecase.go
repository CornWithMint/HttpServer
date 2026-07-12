package usecase

import (
	"errors"
	"server/domain"
	"sync"
	"time"
)

var ErrEmptyTitle = errors.New("Title cannot be empty")

type TaskRepository interface {
	SaveTask(task *domain.Task) (*domain.Task, error)
	GetAllTasks() []*domain.Task
	GetTaskByID(id int) (*domain.Task, error)
}

type TaskUsecase struct {
	taskusecase TaskRepository
	mu          *sync.RWMutex
	nextId      int
}

func NewTaskUsecase(storage TaskRepository) *TaskUsecase {
	return &TaskUsecase{
		taskusecase: storage,
		mu:          &sync.RWMutex{},
		nextId:      0,
	}
}

func (tu *TaskUsecase) CreateTask(title string) (*domain.Task, error) {
	tu.mu.Lock()
	defer tu.mu.Unlock()

	task := &domain.Task{}
	if title == "" {
		return nil, ErrEmptyTitle
	}
	task.ID = tu.nextId
	tu.nextId++
	task.Done = false
	task.CreatedAt = time.Now()
	task, err := tu.taskusecase.SaveTask(task)
	if err != nil {
		return nil, errors.New("Problem With saving")
	}
	return task, nil
}

func (tu *TaskUsecase) ListTasks() ([]*domain.Task, error) {
	arr := tu.taskusecase.GetAllTasks()
	return arr, nil
}

func (tu *TaskUsecase) GetTaskByID(userID int) (*domain.Task, error) {
	task, err := tu.taskusecase.GetTaskByID(userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
