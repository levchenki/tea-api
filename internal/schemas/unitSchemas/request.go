package unitSchemas

import (
	"fmt"
	"github.com/levchenki/tea-api/internal/entity"
	"net/http"
)

type RequestModel struct {
	IsApiece   bool   `json:"isApiece"`
	WeightUnit string `json:"weightUnit"`
	Value      int64  `json:"value"`
}

func (rm *RequestModel) Bind(r *http.Request) error {

	if rm.Value <= 0 {
		return fmt.Errorf("the value must be greater than 0")
	}

	var wu entity.WeightUnit
	err := wu.Scan(rm.WeightUnit)
	if err != nil {
		return fmt.Errorf("invalid weight unit: %s", rm.WeightUnit)
	}
	return nil
}
