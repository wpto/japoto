package expanddb

import (
	"fmt"
	"log"
	"regexp"

	"github.com/pgeowng/japoto/config"
	"github.com/pgeowng/japoto/namematch"
	"github.com/pgeowng/japoto/store"
	"github.com/pgeowng/japoto/types"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expanddb",
		Short: "Parses db and add new info",
		Long:  `expands channel info by adding new parsed values needed later`,
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {

	store := store.NewFileStore()

	entries := store.Read()
	entries = ExtendContent(entries)
	store.Write(entries)
	log.Println("Done")
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

func ExtendContent(eps []types.Entry) []types.Entry {
	for idx := range eps {
		eps[idx].Provider = "unknown"
		eps[idx].Date = "000000"
		eps[idx].ShowId = "unknown"

		info, err := namematch.ExtractInfo(eps[idx].Filename)
		if err != nil {
			log.Fatal(err)
		}

		eps[idx].Date = info.Date
		eps[idx].ShowId = info.ShowId
		eps[idx].Provider = info.Provider

		if len(eps[idx].Title) == 0 {
			eps[idx].Title = fmt.Sprintf("%s %s", eps[idx].Date, eps[idx].ShowId)
		}

		prefixTitleRE := regexp.MustCompile(`^(\d{6})`)
		match := prefixTitleRE.FindStringSubmatch(eps[idx].Title)
		if len(match) == 0 {
			eps[idx].Title = fmt.Sprintf("%s %s", eps[idx].Date, eps[idx].Title)
		}

		eps[idx].URL = config.ChannelPrefix + fmt.Sprint(eps[idx].MessageId)
		seconds := eps[idx].Duration
		minutes := seconds / 60
		eps[idx].DurationHuman = fmt.Sprintf("%d:%02d", minutes, seconds%60)

		eps[idx].SizeHuman = HumanSize(eps[idx].Size)
	}

	return eps
}
