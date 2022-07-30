package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"skhaz.dev/rest/model"
	"skhaz.dev/rest/repository"
)

type WorkspaceRepository struct {
	err        error
	workspace  *model.Workspace
	workspaces model.WorkspaceCollection
	success    bool
}

func (r *WorkspaceRepository) Configure(db *gorm.DB) {
}

func (r *WorkspaceRepository) List(after time.Time, limit int) (any, error) {
	return r.workspaces, r.err
}

func (r *WorkspaceRepository) Get(id any) (any, error) {
	return r.workspace, r.err
}

func (r *WorkspaceRepository) Create(entity any) (any, error) {
	return r.workspace, r.err
}

func (r *WorkspaceRepository) Update(id any, entity any) (bool, error) {
	return r.success, r.err
}

func (r *WorkspaceRepository) Delete(id any) (bool, error) {
	return r.success, r.err
}

func TestGetWorkspacesSuccessfully(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	// workspaces := []model.Workspace{{Name: randstr.String(16)}}
	workspaces := model.WorkspaceCollection{{Name: randstr.String(16)}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{workspaces: workspaces}))

	GetWorkspaces(ctx)

	// namedMap := BuildMultipleResources("workspaces", workspaces).ToMap()
	// payload, _ := json.Marshal(namedMap.Content)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	// assert.Equal(t, payload, r.Body.Bytes())
}

func TestGetWorkspacesAfteWithInvalidTime(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/?after=z", nil)

	GetWorkspaces(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestGetWorkspacesLimitGreaterThanOneHundred(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/?limit=101", nil)

	GetWorkspaces(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestGetWorkspacesError(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{err: gorm.ErrMissingWhereClause}))

	GetWorkspaces(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestGetWorkspaceWithMalformedID(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "0"}}

	GetWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestGetWorkspaceSuccessful(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	uid := "f3edb291-a99d-4a43-8de0-1d6acd00c64d"
	ctx.Params = []gin.Param{{Key: "uuid", Value: uid}}
	workspace := model.Workspace{Name: randstr.String(16)}
	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{workspace: &workspace}))

	GetWorkspace(ctx)

	// namedMap := BuildSingleResource(fmt.Sprintf("/workspaces/%v", uid), workspace).ToMap()
	// payload, _ := json.Marshal(namedMap.Content)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	// assert.Equal(t, payload, r.Body.Bytes())
}

func TestGetWorkspaceError(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "6b7bbb57-179d-47cb-a41d-28d8a429763b"}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{err: gorm.ErrRecordNotFound}))

	GetWorkspace(ctx)

	assert.Equal(t, http.StatusNotFound, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:database/record-not-found")
}

func TestCreateWorkspaceSuccessful(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	workspace := model.Workspace{Name: randstr.String(16)}
	ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer([]byte(fmt.Sprintf(`{"name": "%v"}`, workspace.Name))))

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{workspace: &workspace}))

	CreateWorkspace(ctx)

	assert.Equal(t, http.StatusCreated, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), fmt.Sprintf(`"name":"%v"`, workspace.Name))
}

func TestCreateWorkspaceWithIvalidBody(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"Name":}`)))

	CreateWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestCreateWorkspaceError(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	workspace := model.Workspace{Name: randstr.String(16)}
	ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer([]byte(fmt.Sprintf(`{"name": "%v"}`, workspace.Name))))

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{err: gorm.ErrMissingWhereClause}))

	CreateWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestUpdateWorkspaceSuccessful(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	workspace := model.Workspace{Name: randstr.String(16)}
	ctx.Request, _ = http.NewRequest("PATCH", "/", bytes.NewBuffer([]byte(fmt.Sprintf(`{"name": "%v"}`, workspace.Name))))
	ctx.Params = []gin.Param{{Key: "uuid", Value: "f3edb291-a99d-4a43-8de0-1d6acd00c64d"}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{workspace: &workspace}))

	UpdateWorkspace(ctx)

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Equal(t, "application/hal+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), fmt.Sprintf(`"name":"%v"`, workspace.Name))
}

func TestUpdateWorkspaceWithMalformedBody(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{"Name":}`)))

	UpdateWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestUpdateWorkspaceWithMalformedID(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "0"}}

	UpdateWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestUpdateWorkspaceError(t *testing.T) {

}

func TestDeleteWorkspaceSuccessful(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("DELETE", "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "f3edb291-a99d-4a43-8de0-1d6acd00c64d"}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{}))

	DeleteWorkspace(ctx)

	assert.Equal(t, http.StatusNoContent, r.Code)
	assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
}

func TestDeleteWorkspaceWithMalformedID(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "0"}}

	DeleteWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
	assert.Contains(t, r.Body.String(), "errors:http/bad-request")
}

func TestDeleteWorkspaceError(t *testing.T) {
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request, _ = http.NewRequest("DELETE", "/", nil)
	ctx.Params = []gin.Param{{Key: "uuid", Value: "64675a94-fa11-4544-8f2a-f041b4bdee30"}}

	ctx.Set("RepositoryRegistry", repository.NewRepositoryRegistry(nil, &WorkspaceRepository{err: gorm.ErrMissingWhereClause}))

	DeleteWorkspace(ctx)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "application/problem+json", r.Header().Get("Content-Type"))
}
