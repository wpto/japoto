package printers

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/pgeowng/japoto/config"
	"github.com/pgeowng/japoto/types"
)

type Recent struct{}

func (r Recent) Print(entries []types.Entry) {
	files := []string{
		"./template/base.layout.tmpl",
		"./template/recent.content.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(filepath.Join(config.Dest, "index.html"))
	if err != nil {
		log.Fatalln("index.html create error:", err)
	}

	defer f.Close()

	recent := make(map[string][]types.Entry)
	for _, ep := range entries {
		provider := ep.Provider
		recent[provider] = append(recent[provider], ep)
	}

	for provider, eps := range recent {
		sort.Slice(recent[provider], func(i, j int) bool {
			return recent[provider][i].MessageId > recent[provider][j].MessageId
		})

		filtered := make([]types.Entry, 0)
		nameSet := make(map[string]bool)
		for _, ep := range eps {
			name := ep.ProgramName
			if _, ok := nameSet[name]; !ok {
				filtered = append(filtered, ep)
				nameSet[name] = true
			}
		}

		recent[provider] = filtered

		currLimit := cap(recent[provider])
		if currLimit > config.RecentLimit {
			currLimit = config.RecentLimit
		}
		recent[provider] = recent[provider][:currLimit]
	}

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL": config.PublicURL,
		"Recent":    &recent,
	})

	if err != nil {
		log.Fatalln("index.html write error:", err)
	}
}
