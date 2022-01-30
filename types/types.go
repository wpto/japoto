package types

type Entry struct {
	Size         int    `json:"size"`
	MessageId    int    `json:"message_id"`
	Title        string `json:"title"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer"`
	Filename     string `json:"filename"`
	Date         string `json:"date"`
	ProgramName  string `json:"program_name"`
	Provider     string `json:"provider"`
	URL          string `json:"url"`
	DurationTime string `json:"duration_human"`
	SizeHuman    string `json:"size_human"`
	HasImage     bool   `json:"has_image"`
}

type Source interface {
	GetShows() []Entry
}

type Printer interface {
	Print(db []Entry) error
}
