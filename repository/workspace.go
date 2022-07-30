package repository

import (
	"fmt"
	"time"

	"skhaz.dev/rest/model"
)

type WorkspaceRepository struct {
	GormRepository
}

func (r *WorkspaceRepository) List(after time.Time, limit int) (any, error) {
	var wc model.WorkspaceCollection
	order := "created_at"
	err := r.db.Limit(limit).Order(order).Where(fmt.Sprintf("%v > ?", order), after).Limit(limit).Find(&wc).Error

	return wc, err
}

func (r *WorkspaceRepository) Get(id any) (any, error) {
	var w *model.Workspace

	err := r.db.Where("id = ?", id).First(&w).Error

	return w, err
}

func (r *WorkspaceRepository) Create(entity any) (any, error) {
	w := entity.(*model.Workspace)

	err := r.db.Create(w).Error

	return w, err
}

func (r *WorkspaceRepository) Update(id any, entity any) (bool, error) {
	w := entity.(*model.Workspace)

	if err := r.db.Model(w).Where("id = ?", id).Updates(w).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *WorkspaceRepository) Delete(id any) (bool, error) {
	if err := r.db.Delete(&model.Workspace{}, "id = ?", id).Error; err != nil {
		return false, err
	}

	return true, nil
}
