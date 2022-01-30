package source

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/pgeowng/japoto/types"
)

type FileSource struct {
	Srcpath string
}

func (fs *FileSource) GetShows() []types.Entry {
	data, err := ioutil.ReadFile(fs.Srcpath)
	if err != nil {
		log.Fatalf("error: read %s: %s\n", fs.Srcpath, err)
	}

	result := make([]types.Entry, 0)
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("error: parse %s: %s\n", fs.Srcpath, err)
	}

	return result
}
