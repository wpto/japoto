package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	router := httprouter.New()
	router.GET("/", index)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := sql.Open("sqlite3", "../japoto-dl/japoto-archive.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query("select key, msg_id, size, duration from channel order by msg_id desc")
	if err != nil {
		log.Fatal(err)
	}

	result := make([]ChannelEntry, 0)
	for rows.Next() {
		var entry ChannelEntry
		size := &entry.Size
		duration := &entry.Duration
		err := rows.Scan(&entry.Key, &entry.MsgID, &size, &duration)
		if err != nil {
			log.Fatal(err)
		}
		if size != nil {
			entry.Size = *size
		}
		if duration != nil {
			entry.Duration = *duration
		}

		entry.SizeHuman = FormatSize(entry.Size)
		entry.DurationHuman = FormatDuration(entry.Duration)
		result = append(result, entry)
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	// err = json.NewEncoder(w).Encode(result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// render template and serve
	t, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, result)
	if err != nil {
		log.Fatal(err)
	}

}

type ChannelEntry struct {
	Key           string `json:"key"`
	MsgID         int    `json:"msg_id"`
	Size          int    `json:"size"`
	Duration      int    `json:"duration"`
	SizeHuman     string `json:"size_human"`
	DurationHuman string `json:"duration_human"`
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
