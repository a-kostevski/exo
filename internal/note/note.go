package note

// CoreNote contains the basic note information
type CoreNote interface {
	ID() string
	Title() string
	Content() string
	Path() string
}

// TemplateNote defines template-related operations
type TemplateNote interface {
	TemplateName() string
	SetTemplateName(name string) error
	ApplyTemplate(data interface{}) error
}

// ContentNote defines content management operations
type ContentNote interface {
	SetContent(content string) error
	SetTitle(title string) error
}

// MetadataNote defines metadata operations
type MetadataNote interface {
	Metadata() map[string]interface{}
	SetMetadata(key string, value interface{}) error
}

// PersistentNote defines persistence operations
type PersistentNote interface {
	Save() error
	Load() error
	Delete() error
}

// Note combines all note functionality
type Note interface {
	CoreNote
	TemplateNote
	ContentNote
	MetadataNote
	PersistentNote
}
