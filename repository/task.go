package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
	"fmt"	
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(taskID int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	filebased *filebased.Data
}

func NewTaskRepo(filebasedDb *filebased.Data) *taskRepository {
	return &taskRepository{
		filebased: filebasedDb,
	}
}

func (t *taskRepository) Store(task *model.Task) error {
	t.filebased.StoreTask(*task)

	return nil
}

func (t *taskRepository) Update(taskID int, task *model.Task) error {
	isExists, err := t.filebased.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("error fetching task: %v", err)
	}

	// Perbarui task yang ada dengan data baru
	isExists.Title = task.Title
	isExists.Priority = task.Priority
	isExists.CategoryID = task.CategoryID
	isExists.Deadline = task.Deadline
	isExists.Status = task.Status

	// Simpan task yang telah diperbarui ke database
	err = t.filebased.UpdateTask(taskID, *isExists)
	if err != nil {
		x := fmt.Errorf("error updating task: %v", err)
		return x
	}

	return nil
}

func (t *taskRepository) Delete(id int) error {
	err := t.filebased.DeleteTask(id)
	if err != nil {
		if err.Error() == "record not found" {
			x := fmt.Errorf("record not found")
			return x
		}
		y := fmt.Errorf("error deleting task: %v", err) 
		return y
	}
	return nil
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	task, err := t.filebased.GetTaskByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, fmt.Errorf("record not found") 
		}
		return nil, err 
	}
	return task, nil
}

func (t *taskRepository) GetList() ([]model.Task, error) {
	listTask, err := t.filebased.GetTasks()
	if err != nil {
		x := fmt.Errorf("error fetching tasks: %v", err)
		return nil, x
	}
	return listTask, nil
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	var catTasks []model.TaskCategory
	listCatTask, err := t.filebased.GetTasks()
	if err != nil {
		x := fmt.Errorf("failed to retrieve tasks: %w", err)
		return nil, x
	}

	for _, task := range listCatTask {
		if task.CategoryID == id {
			catList, err := t.filebased.GetCategoryByID(id)
			if err != nil {
				x := fmt.Errorf("failed to retrieve category: %w", err)
				return nil, x
			}
			catTasks = append(catTasks, model.TaskCategory{
				ID:       task.ID,
				Title:    task.Title,
				Category: catList.Name,
			})
		}
	}

	if len(catTasks) == 0 {
		x := fmt.Errorf("no tasks found for category ID: %d", id)
		return nil, x
	}

	return catTasks, nil
}
