package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dsn                    = "file::memory:?cache=shared"
	opts                   = gorm.Config{}
	repositoryName         = "TestRepository"
	nonExistRepositoryName = "NonExistRepository"
)

type TestRepository struct {
}

func (r *TestRepository) Configure(db *gorm.DB) {
}

func (r *TestRepository) List(after time.Time, limit int) (any, error) {
	return nil, nil
}

func (r *TestRepository) Get(id any) (any, error) {
	return nil, nil
}

func (r *TestRepository) Create(entity any) (any, error) {
	return nil, nil
}

func (r *TestRepository) Update(id any, entity any) (bool, error) {
	return true, nil
}

func (r *TestRepository) Delete(id any) (bool, error) {
	return true, nil
}

func TestNewRepositoryRegistry(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(dsn), &opts)
	assert.NoError(t, err)

	registry := NewRepositoryRegistry(
		db,
		&TestRepository{},
	)

	assert.NotNil(t, registry)
}

func TestRepositoryGet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(dsn), &opts)
	assert.NoError(t, err)

	registry := NewRepositoryRegistry(
		db,
		&TestRepository{},
	)

	repository, err := registry.Repository(repositoryName)
	assert.NotNil(t, repository)
	assert.NoError(t, err)
}

func TestRepositoryNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(dsn), &opts)
	assert.NoError(t, err)

	registry := NewRepositoryRegistry(
		db,
		&TestRepository{},
	)

	repository, err := registry.Repository(nonExistRepositoryName)
	assert.Nil(t, repository)
	assert.Error(t, err)
}

func TestMustRepositoryGet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(dsn), &opts)
	assert.NoError(t, err)

	registry := NewRepositoryRegistry(
		db,
		&TestRepository{},
	)

	repository := registry.MustRepository(repositoryName)
	assert.NotNil(t, repository)
}

func TestMustRepositoryNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(dsn), &opts)
	assert.NoError(t, err)

	registry := NewRepositoryRegistry(
		db,
		&TestRepository{},
	)

	assert.Panics(t, func() { registry.MustRepository(nonExistRepositoryName) })
}
