package tagSchemas

import (
	"fmt"
	"net/http"
	"regexp"
)

type RequestModel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (rm *RequestModel) Bind(r *http.Request) error {
	if rm.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	if rm.Color == "" {
		return fmt.Errorf("color is a required field")
	}

	if match, _ := regexp.Match("^#(?:[0-9a-fA-F]{3}){1,2}$", []byte(rm.Color)); !match {
		return fmt.Errorf("invalid color format")
	}

	return nil
}
