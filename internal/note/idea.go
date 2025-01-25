package note

type IdeaStatus string

const (
	StatusNew         IdeaStatus = "new"
	StatusInProgress  IdeaStatus = "in-progress"
	StatusImplemented IdeaStatus = "implemented"
	StatusArchived    IdeaStatus = "archived"
)

type IdeaNote struct {
	*BaseNote
	status   IdeaStatus
	tags     []string
	category string
}

func NewIdeaNote(title string, opts ...NoteOption) (*IdeaNote, error) {
	// Create base note with options
	baseNote, err := NewBaseNote(title, opts...)
	if err != nil {
		return nil, err
	}

	// idea, err := noteService.CreateNote("idea", title,
	// 	note.WithTemplateName("idea"),
	// 	note.WithSubDir(filepath.Join(cfg.Dir.IdeaDir, time.Now().Format("2006/01"))),
	// 	note.WithFileName(fmt.Sprintf("%s.md", title)),
	// )

	// Set default template
	if err := WithTemplateName("idea")(baseNote); err != nil {
		return nil, err
	}

	return &IdeaNote{
		BaseNote: baseNote,
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
