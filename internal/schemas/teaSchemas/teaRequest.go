package teaSchemas

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type RequestModel struct {
	Name        string      `json:"name"`
	Price       float64     `json:"price"`
	Description string      `json:"description"`
	CategoryId  uuid.UUID   `json:"categoryId"`
	TagIds      []uuid.UUID `json:"tagIds,omitempty"`
}

func (tr *RequestModel) Bind(r *http.Request) error {
	if tr.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}
