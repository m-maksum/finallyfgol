package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
	"fmt"
)

type CategoryRepository interface {
	Store(Category *model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	GetByID(id int) (*model.Category, error)
	GetList() ([]model.Category, error)
}

type categoryRepository struct {
	filebasedDb *filebased.Data
}

func NewCategoryRepo(filebasedDb *filebased.Data) *categoryRepository {
	return &categoryRepository{filebasedDb}
}

func (c *categoryRepository) Store(Category *model.Category) error {
	c.filebasedDb.StoreCategory(*Category)
	return nil
}

func (c *categoryRepository) Update(id int, category model.Category) error {
	upCat, err := c.filebasedDb.GetCategoryByID(id)
	if err != nil {
		return fmt.Errorf("error fetching category: %v", err)
	}

	upCat.Name = category.Name

	err = c.filebasedDb.UpdateCategory(id, *upCat)
	if err != nil {
		x := fmt.Errorf("error updating category: %v", err)
		return x
	}

	return nil
}

func (c *categoryRepository) Delete(id int) error {
	err := c.filebasedDb.DeleteCategory(id)
	if err != nil {
		x := fmt.Errorf("error deleting category: %v", err)
		return x
	}
	return nil
}

func (c *categoryRepository) GetByID(id int) (*model.Category, error) {
	idCat, err := c.filebasedDb.GetCategoryByID(id)

	return idCat, err
}

func (c *categoryRepository) GetList() ([]model.Category, error) {
	listCat, err := c.filebasedDb.GetCategories()
	if err != nil {
		x := fmt.Errorf("error fetching categories: %v", err)
		return nil, x
	}
	return listCat, nil}
