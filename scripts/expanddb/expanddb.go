package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pgeowng/japoto/types"

	_ "github.com/mattn/go-sqlite3"
)

// func Cmd() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "expanddb",
// 		Short: "Parses db and add new info",
// 		Long:  `expands channel info by adding new parsed values needed later`,
// 		Run:   run,
// 	}
// }

// func run(cmd *cobra.Command, args []string) {

// 	store := store.NewFileStore(config.FileStorePath)

// 	entries := store.Read()
// 	entries = ExtendContent(entries)
// 	// entries = ExtendPerformers(entries)
// 	store.Write(entries)
// 	log.Println("Done")
// }

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

func main() {
	if err := UpdateTitle(); err != nil {
		log.Fatal(err)
	}
}

func UpdateTitle() (err error) {
	db, err := sql.Open("sqlite3", "scripts/migration/past.db")
	if err != nil {
		return fmt.Errorf("archive open: %w", err)
	}

	rows, err := db.Query("select size, message_id, duration, channel_title, channel_performers, filename, date, show_id from legacy_history where channel_title <> ''")
	if err != nil {
		return fmt.Errorf("archive query channel entry: %w", err)
	}

	defer rows.Close()

	result := make([]ArchiveEntry, 0)
	for rows.Next() {
		var e ArchiveEntry
		err = rows.Scan(&e.Size, &e.MessageID, &e.Duration, &e.Title, &e.Performer, &e.Filename, &e.Date, &e.ShowID)
		if err != nil {
			return fmt.Errorf("archive scan channel entry: %w", err)
		}

		e = ParseTitle(e)

		result = append(result, e)
	}

	_ = result
	for _, e := range result {
		fmt.Println(e)
		_, err = db.Exec("update legacy_history set ep_title = ?, show_title = ? where message_id = ?", e.EpTitle, e.Title, e.MessageID)
		if err != nil {
			return fmt.Errorf("archive update entry: %w", err)
		}
	}

	return
}

func ParseTitle(e ArchiveEntry) ArchiveEntry {
	// fmt.Printf("%#v", e)
	past := e.Title
	e.Title = strings.TrimPrefix(e.Title, e.Date+" ")
	e.Title = strings.TrimPrefix(e.Title, e.ShowID+" ")
	// fmt.Println(e.Title)
	var err error
	/*
		usePartIdx := false
		partIdx := 0

		fmt.Println(e.Title)
		_, err = fmt.Sscanf(e.Title, "p%d", &partIdx)
		if err == nil {
			usePartIdx = true
			e.Title = strings.TrimPrefix(e.Title, "p"+strconv.Itoa(partIdx)+" ")
		}


		_ = usePartIdx
	*/

	epIdx := 0
	_, err = fmt.Sscan(e.Title, &epIdx)
	if err == nil {
		e.Title = strings.TrimPrefix(e.Title, strconv.Itoa(epIdx)+" ")
	}

	_, err = fmt.Sscanf(e.Title, "第%d回", &epIdx)
	if err == nil {
		e.Title = strings.TrimPrefix(e.Title, "第"+strconv.Itoa(epIdx)+"回 ")
		e.EpTitle = "第" + strconv.Itoa(epIdx) + "回"
	}

	if strings.HasPrefix(e.Title, "予告 ") {
		e.EpTitle += " 予告"
		e.Title = strings.TrimPrefix(e.Title, "予告 ")
	}

	if strings.HasPrefix(e.Title, "SP SP回") {
		e.EpTitle += "第SP回"
		e.Title = strings.TrimPrefix(e.Title, "SP SP回 ")
	}

	if strings.HasPrefix(e.Title, "SP 第SP回") {
		e.EpTitle += "第SP回"
		e.Title = strings.TrimPrefix(e.Title, "SP 第SP回 ")
	}

	if strings.HasPrefix(e.Title, "SP2 第SP2回") {
		e.EpTitle += "第SP2回"
		e.Title = strings.TrimPrefix(e.Title, "SP2 第SP2回 ")
	}

	if e.Title != past {
		fmt.Println(e.EpTitle, e.Title, ":", past)
	}

	e.Title = strings.TrimSpace(e.Title)
	e.EpTitle = strings.TrimSpace(e.EpTitle)

	/*
		rr := `(\d{6} )?([^ ]+ )(\d+)? ?(第[^回]+回)(.+)?`
		re := regexp.MustCompile(rr)
		match := re.FindStringSubmatch(e.Title)
		// fmt.Println(match)
		if match != nil {

			date := match[1]
			programID := match[2]
			epIdx := match[3]
			epStr := match[4]
			epTitle := match[5]

			fmt.Printf(
				"%s: (date %s) (program %s) (epIdx %s) (epStr %s) (epTitle %s)\n",
				e.Title,
				date,
				programID,
				epIdx,
				epStr,
				epTitle,
			)
		} else {
			// fmt.Printf("BAD! %s\n", e.Title)
		}
	*/

	// https: //regexr.com/72pvd
	/*
		info, err := ExtendContent(e.Filename)
		if err != nil {
			return fmt.Errorf("archive process: %w", err)
		}
		fmt.Println(info)

		e.Date = info.Date
		e.ShowID = info.ShowId
		e.Source = info.Provider
	*/

	return e
}

func ExtractFromFilename() (err error) {
	db, err := sql.Open("sqlite3", "scripts/migration/past.db")
	if err != nil {
		return fmt.Errorf("archive open: %w", err)
	}

	rows, err := db.Query("select size, message_id, duration, channel_title, channel_performers, filename from legacy_history where filename <> ''")
	if err != nil {
		return fmt.Errorf("archive query channel entry: %w", err)
	}

	defer rows.Close()

	result := make([]ArchiveEntry, 0)
	for rows.Next() {
		var e ArchiveEntry
		err = rows.Scan(&e.Size, &e.MessageID, &e.Duration, &e.Title, &e.Performer, &e.Filename)
		if err != nil {
			return fmt.Errorf("archive scan channel entry: %w", err)
		}

		info, err := ExtendContent(e.Filename)
		if err != nil {
			return fmt.Errorf("archive process: %w", err)
		}
		fmt.Println(info)

		e.Date = info.Date
		e.ShowID = info.ShowId
		e.Source = info.Provider
		result = append(result, e)
	}

	for _, e := range result {
		fmt.Println(e)
		_, err = db.Exec("update legacy_history set date = ?, show_id = ?, source = ? where message_id = ?", e.Date, e.ShowID, e.Source, e.MessageID)
		if err != nil {
			return fmt.Errorf("archive update entry: %w", err)
		}
	}

	return
}

func ExtendContent(filename string) (info EpInfo, err error) {
	// for idx := range eps {
	// eps[idx].Provider = "unknown"
	// eps[idx].Date = "000000"
	// eps[idx].ShowId = "unknown"

	info, err = GuessMeta(filename)
	if err != nil {
		err = fmt.Errorf("extend content: %w", err)
		return
	}

	return

	// eps[idx].Date = info.Date
	// eps[idx].ShowId = info.ShowId
	// eps[idx].Provider = info.Provider

	// if len(eps[idx].Title) == 0 {
	// 	eps[idx].Title = fmt.Sprintf("%s %s", eps[idx].Date, eps[idx].ShowId)
	// }

	// prefixTitleRE := regexp.MustCompile(`^(\d{6})`)
	// match := prefixTitleRE.FindStringSubmatch(eps[idx].Title)
	// if len(match) == 0 {
	// 	eps[idx].Title = fmt.Sprintf("%s %s", eps[idx].Date, eps[idx].Title)
	// }

	// eps[idx].URL = config.ChannelPrefix + fmt.Sprint(eps[idx].MessageId)
	// seconds := eps[idx].Duration
	// minutes := seconds / 60
	// eps[idx].DurationHuman = fmt.Sprintf("%d:%02d", minutes, seconds%60)

	// eps[idx].SizeHuman = entity.HumanSize(eps[idx].Size)

}

func ExtendPerformers(eps []types.Entry) []types.Entry {
	for idx := range eps {
		info, err := GuessPerformers(eps[idx].Performer)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("ok: %v\n", info)
		}
	}
	return eps
}
