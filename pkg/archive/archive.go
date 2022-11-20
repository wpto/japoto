package archive

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pgeowng/japoto/pkg/entity"
)

type Archive struct {
	*sql.DB
}

func NewArchive() (a *Archive, err error) {
	db, err := sql.Open("sqlite3", "scripts/migration/past.db")
	if err != nil {
		err = fmt.Errorf("new archive: %w", err)
		return
	}

	return &Archive{DB: db}, nil
}

func (a *Archive) QueryAllEpisodes() (result []entity.ArchiveEntry, err error) {
	rows, err := a.Query(`select
	source,
	show_id,
	date,
	ep_title,
	onsen_show_id,
	onsen_ep_id,
	onsen_date,
	hibiki_show_id,
	hibiki_ep_id,
	hibiki_video_id,
	size,
	message_id,
	duration,
	channel_title,
	channel_performers,
	filename,
	show_title
	from legacy_history`)
	if err != nil {
		err = fmt.Errorf("archive query all episodes: %w", err)
		return
	}

	defer rows.Close()

	result = make([]entity.ArchiveEntry, 0)
	for rows.Next() {
		var le entity.LegacyArchiveEntry
		err = rows.Scan(
			&le.Source,
			&le.ShowID,
			&le.Date,
			&le.EpTitle,
			&le.OnsenShowID,
			&le.OnsenEpID,
			&le.OnsenDate,
			&le.HibikiShowID,
			&le.HibikiEpID,
			&le.HibikiVideoID,
			&le.Size,
			&le.MessageID,
			&le.Duration,
			&le.ChannelTitle,
			&le.ChannelPerformers,
			&le.Filename,
			&le.ShowTitle,
		)

		if err != nil {
			err = fmt.Errorf("archive query entry: %w", err)
			return
		}

		e := entity.ArchiveEntry{}
		if le.Size != nil {
			e.Size = *le.Size
		}

		if le.MessageID != nil {
			e.MessageID = *le.MessageID
		}

		if le.Duration != nil {
			e.Duration = *le.Duration
		}

		if le.ShowTitle != nil {
			e.Title = *le.ShowTitle
		}

		if le.ChannelPerformers != nil {
			e.Performer = *le.ChannelPerformers
		}

		if le.Filename != nil {
			e.Filename = *le.Filename
		}

		if le.Date != nil {
			e.Date = *le.Date
		}
		e.ShowID = le.ShowID
		e.Source = le.Source

		if le.EpTitle != nil {
			e.EpTitle = *le.EpTitle
		}

		result = append(result, e)
	}

	return
}
