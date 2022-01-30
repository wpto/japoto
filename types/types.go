package types

type Entry struct {
	Size         int    `json:"size"`
	MessageId    int    `json:"message_id"`
	Title        string `json:"title"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer"`
	Filename     string `json:"filename"`
	Date         string
	ProgramName  string
	Provider     string
	URL          string
	DurationTime string
	SizeHuman    string
}

type Source interface {
	GetShows() []Entry
}

type Printer interface {
	Print(db []Entry) error
}
