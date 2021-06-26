package config

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
)

func TestLoadAllowedImports(test *testing.T) {
	for _, data := range []struct {
		name               string
		prepareConfig      func(test *testing.T) (configPath string)
		wantAllowedImports mapset.Set
		wantedErr          assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			configPath := data.prepareConfig(test)

			receivedAllowedImports, receivedErr := LoadAllowedImports(configPath)

			assert.Equal(test, data.wantAllowedImports, receivedAllowedImports)
			data.wantedErr(test, receivedErr)
		})
	}
}
