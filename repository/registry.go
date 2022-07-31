package repository

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Configure(*gorm.DB)
	List(time.Time, int) (any, error)
	Get(any) (any, error)
	Create(any) (any, error)
	Update(any, entity any) (bool, error)
	Delete(any) (bool, error)
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

	return nil, fmt.Errorf("repository %s does not exist", repositoryName)
}

func (r *RepositoryRegistry) MustRepository(repositoryName string) (repository Repository) {
	repository, err := r.Repository(repositoryName)
	if err != nil {
		panic(err.Error())
	}

	return repository
}
