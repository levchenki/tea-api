package categorySchemas

import (
	"fmt"
	"net/http"
)

type RequestModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (tr *RequestModel) Bind(r *http.Request) error {
	if tr.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}
