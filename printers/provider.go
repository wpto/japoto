package printers

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pgeowng/japoto/config"
	"github.com/pgeowng/japoto/types"
)

type Provider struct {
	Name string
}

func (p Provider) Print(entries []types.Entry) {
	files := []string{
		"./template/base.layout.tmpl",
		"./template/provider.content.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return
	}

	fpath := filepath.Join(config.Dest, fmt.Sprintf("%s.html", p.Name))
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatalf("%s create error: %v\n", fpath, err)
	}
	defer f.Close()

	entries = FilterProvider(entries, p.Name)
	entries = UniqueRecentShows(entries)

	alphabet := make(map[string][]types.Entry)
	for _, ep := range entries {
		programName := ep.ProgramName
		if len(programName) == 0 {
			log.Fatalf("programName zero\n %v", ep)
		}

		letter := string(programName[0])
		letter = strings.ToLower(letter)
		if "0" <= letter && letter <= "9" {
			letter = "0"
		}
		alphabet[letter] = append(alphabet[letter], ep)
	}

	for _, eps := range alphabet {
		sort.Slice(eps, func(i, j int) bool {
			return strings.ToLower(eps[i].ProgramName) < strings.ToLower(eps[j].ProgramName)
		})
	}

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL": config.PublicURL,
		"Alphabet":  &alphabet,
	})
	if err != nil {
		log.Fatalf("%s. write error: %v\n", fpath, err)
	}
}
