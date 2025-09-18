package i18n

import (
	"encoding/json"

	"github.com/pelletier/go-toml/v2"
)

type UnmarshalFunc func([]byte, interface{}) error

func getUnmarshalFunc(format i18nFileType) UnmarshalFunc {
	switch format {
	case "toml":
		return toml.Unmarshal
	case "json":
		return json.Unmarshal
	default:
		return json.Unmarshal
	}
}
