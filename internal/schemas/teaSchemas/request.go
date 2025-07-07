package teaSchemas

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type RequestModel struct {
	Name        string      `json:"name"`
	ServePrice  float64     `json:"servePrice"`
	UnitPrice   float64     `json:"unitPrice"`
	UnitId      uuid.UUID   `json:"unitId"`
	Description string      `json:"description,omitempty"`
	CategoryId  uuid.UUID   `json:"categoryId"`
	TagIds      []uuid.UUID `json:"tagIds,omitempty"`
	IsHidden    bool        `json:"isHidden,omitempty"`
}

func (tr *RequestModel) Bind(r *http.Request) error {
	if tr.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	if tr.ServePrice <= 0 {
		return fmt.Errorf("servePrice must be greater than zero")
	}
	if tr.UnitPrice <= 0 {
		return fmt.Errorf("unitPrice must be greater than zero")
	}
	if tr.UnitId == uuid.Nil {
		return fmt.Errorf("unitId is a required field")
	}
	if tr.CategoryId == uuid.Nil {
		return fmt.Errorf("categoryId is a required field")
	}

	return nil
}

type Evaluation struct {
	Rating float64 `json:"rating"`
	Note   string  `json:"note"`
}

func (e *Evaluation) Bind(r *http.Request) error {
	if e.Rating < 1 || e.Rating > 10 {
		return fmt.Errorf("rating should be between 1 and 10 ")
	}
	return nil
}
