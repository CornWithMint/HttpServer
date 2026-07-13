package usecase

import (
	"errors"
	"server/domain"
	"sync"
	"time"

	"github.com/google/uuid"
)

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

func (tu *TaskUsecase) CreateTask(userID uuid.UUID, title string) (*domain.Task, error) {
	tu.mu.Lock()
	defer tu.mu.Unlock()

	task := &domain.Task{}
	if title == "" {
		return nil, BadRequest
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

func (tu *TaskUsecase) ListTasks(userID uuid.UUID) ([]*domain.Task, error) {
	arr := tu.taskusecase.GetAllTasks()
	return arr, nil
}

func (tu *TaskUsecase) GetTaskByID(userID uuid.UUID, ID int) (*domain.Task, error) {
	task, err := tu.taskusecase.GetTaskByID(ID)
	if err != nil {
		return nil, err
	}
	if task.UserId != userID {
		return nil, err
	}
	return task, nil
}
