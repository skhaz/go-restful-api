package controller

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"skhaz.dev/rest/model"
	"skhaz.dev/rest/repository"
)

type query struct {
	After time.Time `form:"after"`
	Limit int       `form:"limit,default=10" binding:"gte=1,lte=100"`
}

type params struct {
	ID string `uri:"uuid" validate:"required,uuid4"`
}

func GetWorkspaceRepository(ctx *gin.Context) repository.Repository {
	return ctx.MustGet("RepositoryRegistry").(*repository.RepositoryRegistry).MustRepository("WorkspaceRepository")
}

func GetWorkspaces(ctx *gin.Context) {
	var q = query{}

	if err := ctx.ShouldBindQuery(&q); err != nil {
		HandleError(err, ctx)
		return
	}

	entities, err := GetWorkspaceRepository(ctx).List(q.After, q.Limit)
	if err != nil {
		HandleError(err, ctx)
		return
	}

	WriteHAL(ctx, http.StatusOK, entities.(model.WorkspaceCollection).ToHAL(ctx.Request.URL.Path, ctx.Request.URL.Query()))
}

func GetWorkspace(ctx *gin.Context) {
	p := params{}

	ctx.ShouldBindUri(&p)

	if err := validate.Struct(p); err != nil {
		HandleError(err, ctx)
		return
	}

	entity, err := GetWorkspaceRepository(ctx).Get(p.ID)
	if err != nil {
		HandleError(err, ctx)
		return
	}

	WriteHAL(ctx, http.StatusOK, entity.(*model.Workspace).ToHAL(ctx.Request.URL.Path))
}

func CreateWorkspace(ctx *gin.Context) {
	body := model.Workspace{}

	if err := ctx.BindJSON(&body); err != nil {
		HandleError(err, ctx)
		return
	}

	entity, err := GetWorkspaceRepository(ctx).Create(&body)
	if err != nil {
		HandleError(err, ctx)
		return
	}

	workspace := entity.(*model.Workspace)
	selfHref, _ := url.JoinPath(ctx.Request.URL.Path, workspace.ID.String())
	WriteHAL(ctx, http.StatusCreated, workspace.ToHAL(selfHref))
}

func UpdateWorkspace(ctx *gin.Context) {
	p := params{}

	ctx.ShouldBindUri(&p)

	if err := validate.Struct(p); err != nil {
		HandleError(err, ctx)
		return
	}

	body := model.Workspace{}

	if err := ctx.BindJSON(&body); err != nil {
		HandleError(err, ctx)
		return
	}

	repository := GetWorkspaceRepository(ctx)

	_, err := repository.Update(p.ID, &body)
	if err != nil {
		HandleError(err, ctx)
		return
	}

	entity, err := repository.Get(p.ID)
	if err != nil {
		HandleError(err, ctx)
		return
	}

	WriteHAL(ctx, http.StatusOK, entity.(*model.Workspace).ToHAL(ctx.Request.URL.Path))
}

func DeleteWorkspace(ctx *gin.Context) {
	p := params{}

	ctx.ShouldBindUri(&p)

	if err := validate.Struct(p); err != nil {
		HandleError(err, ctx)
		return
	}

	_, err := GetWorkspaceRepository(ctx).Delete(p.ID)
	if err != nil {
		HandleError(err, ctx)
		return
	}

	WriteNoContent(ctx)
}
