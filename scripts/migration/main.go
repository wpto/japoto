package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "past.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create table if not exists legacy_history (source text, show_id text, date text, ep_title text)")
	if err != nil {
		log.Fatal(err)
	}

	/*
		_, err = db.Exec("alter table legacy_history add column size int")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column message_id int")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column duration int")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column channel_title text")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column channel_performers text")
		if err != nil {
			log.Fatal(err)
		}
	*/

	bytes, err := os.ReadFile("./japoto.json")
	if err != nil {
		log.Fatal(err)
	}

	type ArchiveEntry struct {
		Size      int    `json:"size"`
		MessageID int    `json:"message_id"`
		Duration  int    `json:"duration"`
		Title     string `json:"title"`
		Performer string `json:"performer"`
		Filename  string `json:"filename"`
	}

	result := []ArchiveEntry{}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range result {
		e.Title = strings.TrimSpace(e.Title)
		e.Performer = strings.TrimSpace(e.Performer)
		e.Filename = strings.TrimSpace(e.Filename)
		row := db.QueryRow("select count(*) from legacy_history where size = ? and message_id = ? and duration = ? and channel_title = ? and channel_performers = ? and filename = ?", e.Size, e.MessageID, e.Duration, e.Title, e.Performer, e.Filename)
		var count int
		err := row.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		if count == 0 {
			_, err = db.Exec("insert into legacy_history (size, message_id, duration, channel_title, channel_performers, filename) values (?, ?, ?, ?, ?, ?)", e.Size, e.MessageID, e.Duration, e.Title, e.Performer, e.Filename)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println(len(result))

	/*
		_, err = db.Exec("alter table legacy_history add column onsen_show_id text")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column onsen_ep_id text")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column onsen_date text")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column hibiki_show_id text")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column hibiki_ep_id text")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("alter table legacy_history add column hibiki_video_id text")
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		// Read csv file
		csvFile, _ := os.Open("../my-podcast-private-data/history_old.tsv")
		reader := csv.NewReader(bufio.NewReader(csvFile))
		reader.Comma = '\t'
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}

			source := strings.TrimSpace(line[0])
			show_id := strings.TrimSpace(line[1])
			src_show_id := strings.TrimSpace(line[2])
			src_ep_id := strings.TrimSpace(line[3])
			ep_title := strings.TrimSpace(line[4])
			anything := strings.TrimSpace(line[5])

			if source == "onsen" {
				AddOnsenEntry(db, source, show_id, src_show_id, src_ep_id, ep_title, anything)
			} else if source == "hibiki" {
				AddHibikiEntry(db, source, show_id, src_show_id, src_ep_id, ep_title, anything)
			} else {
				log.Fatal("Unknown source: " + source)
			}

			/*
				source = strings.TrimSpace(line[0])
				show_id = strings.TrimSpace(line[1])
				date = strings.TrimSpace(line[2])
				ep_title = strings.TrimSpace(line[3])

				row := db.QueryRow("select count(*) from legacy_history where source = ? and show_id = ? and date = ? and ep_title = ?", source, show_id, date, ep_title)
				var count int
				err := row.Scan(&count)
				if err != nil {
					log.Fatal(err)
				}

				if count == 0 {
					_, err = db.Exec("insert into legacy_history (source, show_id, date, ep_title) values (?, ?, ?, ?)", source, show_id, date, ep_title)
					if err != nil {
						log.Fatal(err)
					}
				}
		}
	*/
}

func AddOnsenEntry(db *sql.DB, source string, show_id string, onsen_show_id string, onsen_ep_id string, ep_title string, onsen_date string) {
	row := db.QueryRow("select count(*) from legacy_history where source = ? and show_id = ? and onsen_show_id = ? and onsen_ep_id = ? and ep_title = ? and onsen_date = ?", source, show_id, onsen_show_id, onsen_ep_id, ep_title, onsen_date)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		_, err = db.Exec("insert into legacy_history (source, show_id, onsen_show_id, onsen_ep_id, ep_title, onsen_date) values (?, ?, ?, ?, ?, ?)", source, show_id, onsen_show_id, onsen_ep_id, ep_title, onsen_date)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func AddHibikiEntry(db *sql.DB, source string, show_id string, hibiki_show_id string, hibiki_ep_id string, ep_title string, hibiki_video_id string) {
	row := db.QueryRow("select count(*) from legacy_history where source = ? and show_id = ? and hibiki_show_id = ? and hibiki_ep_id = ? and ep_title = ? and hibiki_video_id = ?", source, show_id, hibiki_show_id, hibiki_ep_id, ep_title, hibiki_video_id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		_, err = db.Exec("insert into legacy_history (source, show_id, hibiki_show_id, hibiki_ep_id, ep_title, hibiki_video_id) values (?, ?, ?, ?, ?, ?)", source, show_id, hibiki_show_id, hibiki_ep_id, ep_title, hibiki_video_id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
