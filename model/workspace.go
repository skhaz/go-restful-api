package model

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pmoule/go2hal/hal"
	"gorm.io/gorm"
)

const (
	After = "after"
	Limit = "limit"
)

type Workspace struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"not null,type:text" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,->" json:"-"`
}

type WorkspaceCollection []*Workspace

func (m *Workspace) ToHAL(selfHref string) (root hal.Resource) {
	root = hal.NewResourceObject()
	root.AddData(m)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: selfHref}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	return
}

func (c WorkspaceCollection) ToHAL(selfHref string, queryString url.Values) (root hal.Resource) {
	type NameOnly struct {
		Name string `json:"name"`
	}

	type Result struct {
		Count   int                 `json:"count"`
		Results WorkspaceCollection `json:"results"`
	}

	root = hal.NewResourceObject()

	selfRel := hal.NewSelfLinkRelation()
	// href := strings.Join([]string{selfHref, queryString.Encode()}, "?")
	selfRel.SetLink(&hal.LinkObject{Href: selfHref})
	root.AddLink(selfRel)

	el, hasLast := Last(c)
	if hasLast {
		after, err := el.CreatedAt.MarshalText()
		if NoError(err) {
			queryString.Set(After, string(after))

			nextRel, _ := hal.NewLinkRelation(NextRelation)
			nextLink := &hal.LinkObject{Href: strings.Join([]string{selfHref, queryString.Encode()}, "?")}
			nextRel.SetLink(nextLink)
			root.AddLink(nextRel)
		}
	}

	var embedded []hal.Resource

	for _, i := range c {
		selfLink, _ := hal.NewLinkObject(fmt.Sprintf("%s/%v", selfHref, i.ID))

		selfRel, _ := hal.NewLinkRelation("self")
		selfRel.SetLink(selfLink)

		resource := hal.NewResourceObject()
		resource.AddLink(selfRel)
		resource.AddData(NameOnly{i.Name})

		embedded = append(embedded, resource)
	}

	workspaces, _ := hal.NewResourceRelation("workspaces")
	workspaces.SetResources(embedded)
	root.AddResource(workspaces)
	root.AddData(Result{len(c), c})

	return
}
