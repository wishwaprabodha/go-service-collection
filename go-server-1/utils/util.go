package utils

import (
	"encoding/json"
	"strings"
)

func JSONDecoder(s string) *json.Decoder {
	d := json.NewDecoder(strings.NewReader(s))
	d.UseNumber()
	return d
}
