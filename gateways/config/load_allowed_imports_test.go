package config

import (
	"io/ioutil"
	"os"
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAllowedImports(test *testing.T) {
	for _, data := range []struct {
		name               string
		prepareConfig      func(test *testing.T) (configPath string)
		wantAllowedImports mapset.Set
		wantedErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success without duplicates",
			prepareConfig: func(test *testing.T) (configPath string) {
				config, err := ioutil.TempFile("", "allowed_imports")
				require.NoError(test, err)
				defer config.Close()

				_, err = config.Write([]byte(`["one", "two"]`))
				require.NoError(test, err)

				return config.Name()
			},
			wantAllowedImports: mapset.NewSet("one", "two"),
			wantedErr:          assert.NoError,
		},
		{
			name: "success with duplicates",
			prepareConfig: func(test *testing.T) (configPath string) {
				config, err := ioutil.TempFile("", "allowed_imports")
				require.NoError(test, err)
				defer config.Close()

				_, err = config.Write([]byte(`["one", "two", "two"]`))
				require.NoError(test, err)

				return config.Name()
			},
			wantAllowedImports: mapset.NewSet("one", "two"),
			wantedErr:          assert.NoError,
		},
		{
			name: "error with unmarshalling",
			prepareConfig: func(test *testing.T) (configPath string) {
				config, err := ioutil.TempFile("", "allowed_imports")
				require.NoError(test, err)
				defer config.Close()

				_, err = config.Write([]byte("incorrect"))
				require.NoError(test, err)

				return config.Name()
			},
			wantAllowedImports: nil,
			wantedErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			configPath := data.prepareConfig(test)
			defer os.Remove(configPath)

			receivedAllowedImports, receivedErr := LoadAllowedImports(configPath)

			assert.Equal(test, data.wantAllowedImports, receivedAllowedImports)
			data.wantedErr(test, receivedErr)
		})
	}
}
