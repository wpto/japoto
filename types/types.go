package types

type ArchiveEntry struct {
	Size      int    `json:"size"`
	MessageID int    `json:"message_id"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Performer string `json:"performer"`
	Filename  string `json:"filename"`
	Date      string `json:"date"`
	ShowID    string `json:"show_id"`
	Source    string `json:"source"`
	EpTitle   string `json:"ep_title"`
	PartIdx   int    `json:"part_idx"`
}

type Entry struct {
	Date          string `json:"date"`
	Duration      int    `json:"duration"`
	DurationHuman string `json:"duration_human"`
	Filename      string `json:"filename"`
	HasImage      bool   `json:"has_image"`
	MessageId     int    `json:"message_id"`
	Performer     string `json:"performer"`
	Provider      string `json:"provider"`
	ShowId        string `json:"show_id"`
	Size          int    `json:"size"`
	SizeHuman     string `json:"size_human"`
	Title         string `json:"title"`
	URL           string `json:"url"`
}

// Performers    []Person `json:"performers"`
type Person struct {
	IsGuest   bool    `json:"is_guest"`
	Name      string  `json:"name"`
	Character *string `json:"character"`
}

type Source interface {
	Read() []Entry
	Write([]Entry)
}

type Printer interface {
	Print(db []Entry) error
}
