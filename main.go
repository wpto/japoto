package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pgeowng/japoto/namematch"
)

type Entry struct {
	Size         int    `json:"size"`
	MessageId    int    `json:"message_id"`
	Title        string `json:"title"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer"`
	Filename     string `json:"filename"`
	Date         string
	ProgramName  string
	Provider     string
	URL          string
	DurationTime string
	SizeHuman    string
}

const channelPrefix = "https://t.me/japoto/"
const staticPrefix = "./static"
const outputPrefix = "./public"
const inputFile = "../japoto-private/japoto.json"

// const publicURL = "https://pgeowng.github.io/japoto"

const publicURL = ""

func main() {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("error: read japoto.json: %s\n", err)
	}

	s := make([]Entry, 0)
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("error: parse japoto.json: %s\n", err)
	}

	for idx := range s {
		s[idx].Provider = "unknown"
		s[idx].Date = "000000"
		s[idx].ProgramName = "unknown"

		info, err := namematch.ExtractInfo(s[idx].Filename)
		if err != nil {
			log.Fatal(err)
		}

		s[idx].Date = info.Date
		s[idx].ProgramName = info.ProgramName
		s[idx].Provider = info.Provider

		s[idx].URL = channelPrefix + fmt.Sprint(s[idx].MessageId)
		seconds := s[idx].Duration
		minutes := seconds / 60
		s[idx].DurationTime = fmt.Sprintf("%d:%02d", minutes, seconds%60)

		s[idx].SizeHuman = HumanSize(s[idx].Size)
	}

	db := make(map[string]map[string][]Entry)
	for _, entry := range s {

		_, ok := db[entry.Provider]
		if !ok {
			db[entry.Provider] = make(map[string][]Entry)
		}

		db[entry.Provider][entry.ProgramName] = append(db[entry.Provider][entry.ProgramName], entry)
	}

	err = os.MkdirAll(outputPrefix, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	staticFiles, err := os.ReadDir(staticPrefix)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range staticFiles {
		fp := filepath.Join(staticPrefix, file.Name())
		ff, err := ioutil.ReadFile(fp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dp := filepath.Join(outputPrefix, file.Name())
		err = ioutil.WriteFile(dp, ff, 0644)
		if err != nil {
			fmt.Println("Error creating", dp)
			fmt.Println(err)
			return
		}
	}

	renderIndex(db)
	renderAll(db)

	// sc := presenters.ShowContent(s)
	// presenters.RenderShowContent(sc)
	for provider := range db {
		for showName := range db[provider] {
			renderPage(provider, showName, db[provider][showName])
		}
	}
}

func renderAll(db map[string]map[string][]Entry) {
	files := []string{
		"./template/base.layout.tmpl",
		"./template/all.content.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(filepath.Join(outputPrefix, "all.html"))
	if err != nil {
		log.Fatalln("index.html create error:", err)
	}

	defer f.Close()

	alphabet := make(map[string]map[string][]string)
	for provider, shows := range db {
		for showName := range shows {
			// text, err := json.MarshalIndent(shows[showName], "  ", "  ")
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Printf(string(text))
			letter := string(showName[0])
			letter = strings.ToLower(letter)
			if "0" <= letter && letter <= "9" {
				letter = "numbers"
			}
			if _, ok := alphabet[provider]; !ok {
				alphabet[provider] = make(map[string][]string)
			}
			alphabet[provider][letter] = append(alphabet[provider][letter], showName)
		}
	}

	for _, letters := range alphabet {
		for _, shows := range letters {
			sort.Slice(shows, func(i, j int) bool {
				return strings.ToLower(shows[i]) < strings.ToLower(shows[j])
			})
		}
	}

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL": publicURL,
		"Alphabet":  &alphabet,
	})
	if err != nil {
		log.Fatalln("index.html write error:", err)
	}
}

func renderIndex(db map[string]map[string][]Entry) {
	files := []string{
		"./template/base.layout.tmpl",
		"./template/recent.content.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(filepath.Join(outputPrefix, "index.html"))
	if err != nil {
		log.Fatalln("index.html create error:", err)
	}

	defer f.Close()

	limit := 30
	recent := make(map[string][]Entry)
	for provider, shows := range db {
		for _, eps := range shows {
			recent[provider] = append(recent[provider], eps[0])
		}
		sort.Slice(recent[provider], func(i, j int) bool {
			return recent[provider][i].MessageId > recent[provider][j].MessageId
		})

		currLimit := cap(recent[provider])
		if currLimit > limit {
			currLimit = limit
		}
		recent[provider] = recent[provider][:currLimit]
	}

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL": publicURL,
		"Db":        &db,
		"Recent":    &recent,
	})
	if err != nil {
		log.Fatalln("index.html write error:", err)
	}
}

func renderPage(provider, showName string, entries []Entry) {
	files := []string{
		"./template/base.layout.tmpl",
		"./template/show.content.tmpl",
	}

	err := os.MkdirAll(filepath.Join(outputPrefix, provider), fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	fpath := filepath.Join(outputPrefix, provider, showName+".html")
	f, err := os.Create(fpath)
	if err != nil {
		log.Fatalln(fpath+" create error:", err)
	}

	defer f.Close()

	err = ts.Execute(f, map[string]interface{}{
		"PublicURL": publicURL,
		"Provider":  provider,
		"ShowName":  showName,
		"Entries":   entries,
	})

	if err != nil {
		log.Fatalln(fpath+" write error:", err)
	}
}

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

func TryFilenameV5(filename string, entry *Entry) (ok bool) {
	filenameRE := regexp.MustCompile(`(\d{6})-(.+?)--(.+?).mp3`)
	match := filenameRE.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}

	entry.Date = match[1]
	entry.ProgramName = match[2]
	tags := match[3]

	entry.Provider = "unknown-v5"

	if strings.Contains(tags, "onsen") {
		entry.Provider = "onsen"
	}

	if strings.Contains(tags, "hibiki") {
		entry.Provider = "hibiki"
	}

	return true
}

func TryFilenameV4(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(\d{6})-(.+?)-(\d+?|SP\d*?)-(onsen|hibiki)(-p\d+?)?\.mp3`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	entry.Date = match[1]
	entry.ProgramName = match[2]
	entry.Provider = match[4]
	return true
}

func TryFilenameV3(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(\d{6})-(\d*?)-(.+?)-(onsen|hibiki)(-p\d+?)?\.mp3`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	entry.Date = match[1]
	entry.ProgramName = match[3]
	entry.Provider = match[4]
	return true
}

func TryFilenameV2(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(\d{6})-(.+?)-(\d*?|SP\d*?)?\.mp3`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	entry.Date = match[1]
	entry.ProgramName = match[2]
	entry.Provider = "onsen"
	return true
}

func TryFilenameV1(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(\d{3})(.+?)(\d{6}?)(.{4})\.mp3`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	entry.Date = match[3]
	entry.ProgramName = match[2]
	entry.Provider = "onsen"
	return true
}

func TryFilename100man(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(210508)(-100ma)?`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	entry.Date = "210508"
	entry.ProgramName = "100man"
	entry.Provider = "onsen"
	return true
}
func TryFilenameRadista(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(\d{3})radista_ex_(\d{2})`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	entry.Date = "0000" + match[2]
	entry.ProgramName = "radista_ex"
	entry.Provider = "onsen"
	return true
}
func TryFilenamePhyChe(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`(\d)_(\d+)生肉_PsyChe`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	month, err := strconv.Atoi(match[1])
	if err != nil {
		log.Fatalf("parse error %v", match[1])
	}
	day, err := strconv.Atoi(match[2])
	if err != nil {
		log.Fatalf("parse error %v", match[2])
	}

	entry.Date = fmt.Sprintf("20%02d%02d", month, day)
	entry.ProgramName = "watahana"
	entry.Provider = "onsen"
	return true
}

func TryFilenameAsacoco(filename string, entry *Entry) (ok bool) {
	re := regexp.MustCompile(`【_桐生ココ】あさココ(?:LIVE|ライブ)(?:100回目)?(?:ニュース！|NEWS初回放送)(\d{1,2})\D(\d{1,2})`)
	match := re.FindStringSubmatch(filename)
	if len(match) == 0 {
		return false
	}
	month, err := strconv.Atoi(match[1])
	if err != nil {
		log.Fatalf("parse error %v", match[1])
	}
	day, err := strconv.Atoi(match[2])
	if err != nil {
		log.Fatalf("parse error %v", match[2])
	}

	year := 20
	if month > 11 {
		year = 19
	}

	entry.Date = fmt.Sprintf("%02d%02d%02d", year, month, day)
	entry.ProgramName = "asacoco"
	entry.Provider = "youtube"
	return true
}
