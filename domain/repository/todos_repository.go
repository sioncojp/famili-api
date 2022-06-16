package repository

import (
	"github.com/sioncojp/famili-api/domain"
	"github.com/sioncojp/famili-api/domain/model"
)

// interfaceを使うことでDIPを解決する。mockも作成できるようになる
type TodoRepository interface {
	GetById(domain.Id) (model.Todo, error)
	List() ([]model.Todo, error)
	Create(*model.Todo) error
	Update(*model.Todo) error
	Delete(*model.Todo) error
}
