package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pgeowng/japoto/config"
	"github.com/pgeowng/japoto/printers"
	"github.com/pgeowng/japoto/source"
)

// const channelPrefix = "https://t.me/japoto/"
// const config.Static = "./static"
// const outputPrefix = "./public"
// const inputFile = "../japoto-private/japoto.json"

// const publicURL = "https://pgeowng.github.io/japoto"

const publicURL = ""

func main() {
	ff := source.FileSource{}
	entries := ff.GetShows()

	err := os.MkdirAll(config.Dest, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	staticFiles, err := os.ReadDir(config.Static)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range staticFiles {
		fp := filepath.Join(config.Static, file.Name())
		ff, err := ioutil.ReadFile(fp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dp := filepath.Join(config.Dest, file.Name())
		err = ioutil.WriteFile(dp, ff, 0644)
		if err != nil {
			fmt.Println("Error creating", dp)
			fmt.Println(err)
			return
		}
	}

	arranged := make(map[string]map[string]bool)
	for _, ep := range entries {
		if _, ok := arranged[ep.Provider]; !ok {
			arranged[ep.Provider] = make(map[string]bool)
		}
		arranged[ep.Provider][ep.ProgramName] = true
	}

	r := printers.Recent{}
	r.Print(entries)

	for provider := range arranged {
		pr := printers.Provider{Name: provider}
		pr.Print(entries)
	}

	for provider := range arranged {
		for name := range arranged[provider] {
			sh := printers.Show{
				Provider: provider,
				Name:     name,
			}
			sh.Print(entries)
		}
	}

	// renderIndex(db)
	// renderAll(db)

	// // sc := presenters.ShowContent(s)
	// // presenters.RenderShowContent(sc)
	// for provider := range db {
	// 	for showName := range db[provider] {
	// 		renderPage(provider, showName, db[provider][showName])
	// 	}
	// }
}
