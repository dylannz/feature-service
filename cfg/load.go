package cfg

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func LoadYAML(r io.Reader) (Config, error) {
	dec := yaml.NewDecoder(r)
	cfg := Config{}
	err := dec.Decode(&cfg)
	return cfg, errors.Wrap(err, "read yaml")
}

func LoadYAMLDir(filePath string) (Config, error) {
	cfg := Config{}
	err := filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".yaml", ".yml":
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			c, err := LoadYAML(f)
			if err != nil {
				return err
			}
			cfg.Append(c)
		}

		return nil
	})
	return cfg, err
}
