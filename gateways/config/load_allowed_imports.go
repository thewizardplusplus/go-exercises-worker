package config

import (
	"encoding/json"
	"io/ioutil"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
)

// LoadAllowedImports ...
func LoadAllowedImports(configPath string) (mapset.Set, error) {
	allowedImportsAsJSON, err := ioutil.ReadFile(configPath) // nolint: gosec
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the allowed imports")
	}

	var allowedImportAsSlice []string
	err = json.Unmarshal(allowedImportsAsJSON, &allowedImportAsSlice)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal the allowed imports")
	}

	allowedImports := mapset.NewSet()
	for _, allowedImport := range allowedImportAsSlice {
		allowedImports.Add(allowedImport)
	}

	return allowedImports, nil
}
