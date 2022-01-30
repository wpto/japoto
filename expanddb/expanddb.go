package expanddb

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/pgeowng/japoto/config"
	"github.com/pgeowng/japoto/source"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	src := source.ChannelSource{
		Srcpath: config.ChannelInfoPath,
	}
	entries := src.GetShows()

	text, err := json.Marshal(entries)
	if err != nil {
		log.Fatalf("error when marshaling: %v\n", err)
	}
	err = ioutil.WriteFile(config.FileSourcePath, text, 0644)
	if err != nil {
		log.Fatalf("error when writing: %v\n", err)
	}
	log.Println("Done")
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expanddb",
		Short: "Parses db and add new info",
		Long:  `expands channel info by adding new parsed values needed later`,
		Run:   run,
	}
}
