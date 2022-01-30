package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pgeowng/japoto/config"
	"github.com/pgeowng/japoto/namematch"
	"github.com/pgeowng/japoto/types"
)

type ChannelSource struct {
	Srcpath string
}

func (cs *ChannelSource) GetShows() []types.Entry {

	data, err := ioutil.ReadFile(cs.Srcpath)
	if err != nil {
		log.Fatalf("error: read %s: %s\n", cs.Srcpath, err)
	}

	s := make([]types.Entry, 0)
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("error: parse %s: %s\n", cs.Srcpath, err)
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

		s[idx].URL = config.ChannelPrefix + fmt.Sprint(s[idx].MessageId)
		seconds := s[idx].Duration
		minutes := seconds / 60
		s[idx].DurationTime = fmt.Sprintf("%d:%02d", minutes, seconds%60)

		s[idx].SizeHuman = HumanSize(s[idx].Size)
	}
	return s
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
