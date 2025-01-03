package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var UserConfig, _ = os.UserConfigDir()

var ColorGenConfig = fmt.Sprintf("%v/colorgen", UserConfig)

type Style struct {
	Theme string `yaml:"theme"`

	StyleName string `yaml:"name"`

	Bold  bool `yaml:"bold"`
	Under bool `yaml:"under"`
	Blink bool `yaml:"blink"`

	ForeColor int `yaml:"fore"`
	BackColor int `yaml:"back"`

	FileTypes []string `yaml:"filetypes"`
}

func LoadStyle(theme, styleName string) *Style {
	file, err := os.Open(filepath.Join(ColorGenConfig, theme, styleName+".yaml"))
	if err != nil {
		return nil
	}
	defer file.Close()

	var style *Style
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&style); err != nil {
		return nil
	}

	return style
}

func SaveStyle(style Style) error {
	path := filepath.Join(ColorGenConfig, style.Theme)
	file, err := os.Create(filepath.Join(path, style.StyleName+".yaml"))
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(style)
}
