package database

import (
	"gorm.io/gorm"

	"github.com/sioncojp/famili-api/domain"
	"github.com/sioncojp/famili-api/domain/model"
	"github.com/sioncojp/famili-api/domain/repository"
)

// todoRepository...
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepositoryMySQL...Repository interfaceを返すことでserviceとメソッドを揃える
func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{db}
}

// GetById...IDからtodoを取得するためのDB操作
func (r *todoRepository) GetById(id domain.Id) (model.Todo, error) {
	var result model.Todo
	if err := r.db.Where("id = ?", id).First(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

// List...todoを全て取得するためのDB操作
func (r *todoRepository) List() ([]model.Todo, error) {
	var result []model.Todo

	if err := r.db.Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}

// Create...todo作成するためのDB操作
func (r *todoRepository) Create(todo *model.Todo) error {
	return r.db.Create(&todo).Error
}

// Update...todo更新するためのDB操作
func (r *todoRepository) Update(todo *model.Todo) error {
	return r.db.Save(&todo).Error
}

// Delete...IDからtodo削除するためのDB操作
func (r *todoRepository) Delete(todo *model.Todo) error {
	return r.db.Delete(&model.Todo{}, todo.ID).Error
}
