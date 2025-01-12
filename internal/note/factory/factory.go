package factory

import (
	"time"

	"github.com/a-kostevski/exo/internal/note/api"
	"github.com/a-kostevski/exo/internal/note/types"
	"github.com/a-kostevski/exo/internal/templates"
)

type NoteFactory struct {
	tm templates.TemplateManager
}

func New(tm templates.TemplateManager) *NoteFactory {
	return &NoteFactory{tm: tm}
}

func (f *NoteFactory) CreateDaily(date time.Time) (api.PeriodNote, error) {
	return types.NewDailyNote(date, f.tm)
}

func (f *NoteFactory) CreateZettel(title string) (api.Note, error) {
	return types.NewZettelNote(title, f.tm)
}
