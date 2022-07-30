package repository

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Configure(*gorm.DB)
	List(after time.Time, limit int) (any, error)
	Get(id any) (any, error)
	Create(entity any) (any, error)
	Update(id any, entity any) (bool, error)
	Delete(id any) (bool, error)
}

type GormRepository struct {
	db *gorm.DB
}

func (r *GormRepository) Configure(db *gorm.DB) {
	r.db = db
}

type RepositoryRegistry struct {
	registry map[string]Repository

	db *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB, v ...Repository) *RepositoryRegistry {
	r := &RepositoryRegistry{
		db:       db,
		registry: map[string]Repository{},
	}

	r.registerRepositories(v)
	return r
}

func (r *RepositoryRegistry) registerRepositories(repositories []Repository) {
	for _, v := range repositories {
		repositoryName := reflect.TypeOf(v).Elem().Name()
		v.Configure(r.db)
		r.registry[repositoryName] = v
	}
}

func (r *RepositoryRegistry) Repository(repositoryName string) (Repository, error) {
	if repository, ok := r.registry[repositoryName]; ok {
		return repository, nil
	}

	return nil, errors.New(fmt.Sprintf("repository %s does not exist", repositoryName))
}

func (r *RepositoryRegistry) MustRepository(repositoryName string) (repository Repository) {
	repository, err := r.Repository(repositoryName)
	if err != nil {
		panic(err.Error())
	}

	return repository
}
