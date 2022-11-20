package entity

import "fmt"

func HumanSize(intSize int) string {
	unit := "B"
	size := float64(intSize)
	if size*10 >= 1024 {
		unit = "KB"
		size = size / 1024
	}

	if size*10 >= 1024 {
		unit = "MB"
		size = size / 1024
	}

	if size < 10 {
		return fmt.Sprintf("%.1f%s", size, unit)
	} else {
		return fmt.Sprintf("%.f%s", size, unit)
	}
}

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

type LegacyArchiveEntry struct {
	Source            string
	ShowID            string
	Date              *string
	EpTitle           *string
	OnsenShowID       *string
	OnsenEpID         *string
	OnsenDate         *string
	HibikiShowID      *string
	HibikiEpID        *string
	HibikiVideoID     *string
	Size              *int
	MessageID         *int
	Duration          *int
	ChannelTitle      *string
	ChannelPerformers *string
	Filename          *string
	ShowTitle         *string
}

func FormatSize(size int) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KiB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MiB", float64(size)/1024/1024)
	} else {
		return fmt.Sprintf("%.2f GiB", float64(size)/1024/1024/1024)
	}
}

func FormatDuration(duration int) string {
	if duration < 60*60 {
		return fmt.Sprintf("%02d:%02d", duration/60, duration%60)
	} else {
		return fmt.Sprintf("%d:%02d:%02d", duration/60/60, duration/60%60, duration%60)
	}
}

type EpisodeRender struct {
	Date          string `json:"date"`
	Duration      string `json:"duration"`
	Filename      string `json:"filename"`
	HasImage      bool   `json:"has_image"`
	MessageID     int    `json:"message_id"`
	HasPerformer  bool
	Performer     string `json:"performer"`
	Source        string `json:"provider"`
	ShowID        string `json:"show_id"`
	Size          string `json:"size"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	HasChannelURL bool
	ChannelURL    string `json:"channel_url"`

	EpTitle string
}

func (a ArchiveEntry) Render(channelFmt string) (b EpisodeRender) {
	b.HasChannelURL = a.MessageID != 0
	if b.HasChannelURL {
		b.ChannelURL = fmt.Sprintf(channelFmt, a.MessageID)
	}
	b.HasPerformer = a.Performer != ""
	b.Performer = a.Performer
	b.MessageID = a.MessageID
	b.Filename = a.Filename
	b.Duration = FormatDuration(a.Duration)
	b.Size = FormatSize(a.Size)

	b.Date = a.Date
	b.Source = a.Source
	b.ShowID = a.ShowID
	b.EpTitle = a.EpTitle

	b.Title = fmt.Sprintf("%s %s %s %s %s", a.Date, a.Source, a.ShowID, a.EpTitle, a.Title)

	return
}
