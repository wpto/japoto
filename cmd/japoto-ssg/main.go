package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pgeowng/japoto/pkg/archive"
	"github.com/pgeowng/japoto/pkg/entity"

	_ "github.com/mattn/go-sqlite3"
)

var arch *archive.Archive

func main() {

	a, err := archive.NewArchive()
	if err != nil {
		log.Fatal(err)
	}

	arch = a

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/recent", ShowHandler)
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	entries, err := arch.QueryAllEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	result := make([]ChannelEntry, 0)
	for _, e := range entries {

		/*
			var entry ChannelEntry
			size := &entry.Size
			duration := &entry.Duration

				entry.Size = *size

				entry.Duration = *duration

			entry.SizeHuman = FormatSize(entry.Size)
			entry.DurationHuman = FormatDuration(entry.Duration)
		*/
		entry := ChannelEntry{
			Key:           e.ShowID,
			MsgID:         e.MessageID,
			Size:          e.Size,
			Duration:      e.Duration,
			SizeHuman:     entity.FormatSize(e.Size),
			DurationHuman: entity.FormatDuration(e.Duration),
		}
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

func ShowHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ts, err := template.ParseFiles("template/base.layout.html", "template/show.content.html")
	if err != nil {
		log.Fatal(err)
	}

	entries, err := arch.QueryAllEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	//entries = entries[8000:8005]

	sort.Slice(entries, func(i, j int) bool {
		return entries[j].MessageID <= entries[i].MessageID &&
			entries[j].ShowID <= entries[i].ShowID &&
			entries[j].Date <= entries[i].Date
	})

	renderData := make([]entity.EpisodeRender, len(entries))
	for i, e := range entries {
		fmt.Println(e)
		renderData[i] = e.Render("https://t.me/japoto/%d")
	}

	fmt.Println(renderData)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	err = ts.Execute(w, map[string]interface{}{
		"PublicURL":  ".",
		"CreateTime": time.Now(),
		"Provider":   "all",
		"ShowName":   "all",
		"Episodes":   renderData,
	})
}

type ChannelEntry struct {
	Key           string `json:"key"`
	MsgID         int    `json:"msg_id"`
	Size          int    `json:"size"`
	Duration      int    `json:"duration"`
	SizeHuman     string `json:"size_human"`
	DurationHuman string `json:"duration_human"`
}
