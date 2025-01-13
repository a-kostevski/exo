package builtin

import (
	"github.com/a-kostevski/exo/internal/note"
	"github.com/a-kostevski/exo/internal/templates"
)

type IdeaStatus string

const (
	StatusNew         IdeaStatus = "new"
	StatusInProgress  IdeaStatus = "in-progress"
	StatusImplemented IdeaStatus = "implemented"
	StatusArchived    IdeaStatus = "archived"
)

type IdeaNote struct {
	*note.BaseNote
	status   IdeaStatus
	tags     []string
	category string
}

func NewIdeaNote(title string, tm templates.TemplateManager) (*IdeaNote, error) {
	base, err := note.NewBaseNote(title, tm)
	if err != nil {
		return nil, err
	}

	return &IdeaNote{
		BaseNote: base,
		status:   StatusNew,
		tags:     []string{},
	}, nil
}

func (n *IdeaNote) Status() IdeaStatus { return n.status }
func (n *IdeaNote) Tags() []string     { return n.tags }
func (n *IdeaNote) Category() string   { return n.category }

func (n *IdeaNote) SetStatus(status IdeaStatus) error {
	n.status = status
	// n.base.modified = time.Now()
	return nil
}

func (n *IdeaNote) AddTag(tag string) error {
	n.tags = append(n.tags, tag)
	// n.modified = time.Now()
	return nil
}

func (n *IdeaNote) SetCategory(category string) error {
	n.category = category
	// n.modified = time.Now()
	return nil
}
