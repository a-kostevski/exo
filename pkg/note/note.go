package note

import "time"

type NoteContentManager interface {
	Content() string
	SetContent(content string) error
}

type NoteFS interface {
	Save() error
	Load() error
	Delete() error
}

type Editor interface {
	Open() error
}

type Metadata interface {
	Title() string
	Path() string
	Created() time.Time
	Modified() time.Time
}

type NoteValidator interface {
	Validate() error
}

type Note interface {
	NoteContentManager
	NoteFS
	Editor
	Metadata
	NoteValidator
	String() string
}
