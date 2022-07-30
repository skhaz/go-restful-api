package model

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func TestWorkspaceJsonMarshal(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	name := randstr.String(16)

	w1 := struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w2 := Workspace{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	b1, _ := json.Marshal(w1)
	b2, err := json.Marshal(w2)
	if assert.NoError(t, err) {
		assert.Equal(t, string(b1), string(b2))
	}
}

func TestSingleWorkspaceHAL(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	name := randstr.String(16)
	url := "/" + randstr.String(16)

	w := Workspace{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type Self struct {
		Href string `json:"href"`
	}

	type Links struct {
		Self Self `json:"self"`
	}

	type HAL struct {
		Links     Links     `json:"_links"`
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	expected, _ := json.Marshal(HAL{Links: Links{
		Self: Self{
			Href: url,
		},
	},
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	})

	resource := w.ToHAL(url)
	namedMap := resource.ToMap()
	actual, _ := json.Marshal(namedMap.Content)
	assert.Equal(t, string(expected), string(actual))
}

func TestMultipleWorkspacesHAL(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	name := randstr.String(16)
	path := "/" + randstr.String(16)
	queryString := url.Values{}

	w := Workspace{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	wc := WorkspaceCollection{}
	wc = append(wc, &w)

	type Href struct {
		Href string `json:"href"`
	}

	type Link struct {
		Self Href `json:"self"`
	}

	type LinkWithNext struct {
		Next Href `json:"next"`
		Self Href `json:"self"`
	}

	type Links struct {
		Link Link   `json:"_links"`
		Name string `json:"name,omitempty"`
	}

	type Embedded struct {
		Workspaces []Links `json:"workspaces"`
	}

	type HAL struct {
		Embedded Embedded            `json:"_embedded"`
		Link     LinkWithNext        `json:"_links"`
		Count    int                 `json:"count"`
		Results  WorkspaceCollection `json:"results"`
	}

	after, _ := w.CreatedAt.MarshalText()
	query := url.Values{}
	query.Add(After, string(after))

	hal, _ := json.Marshal(HAL{
		Embedded: Embedded{
			Workspaces: []Links{
				{
					Link{
						Self: Href{
							path + "/" + id.String(),
						},
					},
					name,
				},
			}},
		Link: LinkWithNext{
			Next: Href{strings.Join([]string{path, query.Encode()}, "?")},
			Self: Href{path},
		},
		Count:   1,
		Results: wc,
	})

	resource := wc.ToHAL(path, queryString)
	namedMap := resource.ToMap()
	actual, _ := json.Marshal(namedMap.Content)
	expected, _ := JSONRemarshal(hal) // Sort keys to match with HAL's marshaling
	assert.Equal(t, string(expected), string(actual))
}
