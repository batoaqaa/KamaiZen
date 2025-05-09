package document_manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"iter"
	"maps"
)

//go:embed cookbooks/cookbook_devel.json
var cookbooksFile []byte

type DocEntry struct {
	Name          string `json:"name"`
	Documentation string `json:"documentation"`
}

type Docs struct {
	Docs []DocEntry `json:"docs"`
}

var CookBookDocs map[string]string

func readJSONFromFile() error {
	if len(cookbooksFile) == 0 {
		log.Error().Msg("Cookbooks file is empty")
		return fmt.Errorf("cookbooks file is empty")
	}
	var docs Docs
	err := json.Unmarshal(cookbooksFile, &docs)
	if err != nil {
		log.Error().Err(err).Msg("Error reading JSON")
		return err
	}
	for _, doc := range docs.Docs {
		CookBookDocs[doc.Name] = doc.Documentation
	}
	return nil

}

func init() {
	log.Debug().Msg("Initializing CookBookDocs")
	CookBookDocs = make(map[string]string)
	readJSONFromFile()
}

func GetCookBookDocs(name string) string {
	return CookBookDocs[name]
}

func GetAllCookBookKeys() iter.Seq[string] {
	return maps.Keys(CookBookDocs)
}
