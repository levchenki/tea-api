package teaSchemas

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sort"
	"strconv"
)

type Filters struct {
	Limit      uint64       `json:"limit" db:"limit"`
	Offset     uint64       `json:"offset" db:"offset"`
	CategoryId uuid.UUID    `json:"categoryId,omitempty" db:"category_id"`
	Name       string       `json:"name,omitempty" db:"name"`
	Tags       []string     `json:"tags,omitempty" db:"tags"`
	MinPrice   float64      `json:"minPrice,omitempty" db:"min_price"`
	MaxPrice   float64      `json:"maxPrice,omitempty" db:"max_price"`
	SortBy     SortByFilter `json:"sortBy,omitempty"`
	IsAsc      bool         `json:"isAsc"`
	UserId     uuid.UUID    `json:"userId" db:"user_id"`
}

type SortByFilter string

const (
	NAME  SortByFilter = "name"
	PRICE SortByFilter = "price"
)

func (f *SortByFilter) String() string {
	return string(*f)
}
func (f *SortByFilter) Parse(s string) error {
	teaSortByFilterMap := map[string]SortByFilter{
		"name":  NAME,
		"price": PRICE,
	}
	if val, ok := teaSortByFilterMap[s]; ok {
		*f = val
		return nil
	}
	return fmt.Errorf("invalid TeaSortByFilter value: %s", s)
}

func (tf *Filters) Validate(r *http.Request) error {
	query := r.URL.Query()
	limit, err := strconv.ParseUint(query.Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseUint(query.Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}

	tf.Limit = limit
	tf.Offset = offset

	categoryIdStr := query.Get("categoryId")
	if categoryIdStr != "" {
		categoryId, err := uuid.Parse(categoryIdStr)
		if err != nil {
			categoryId = uuid.Nil
		}
		tf.CategoryId = categoryId
	}

	nameStr := query.Get("name")
	if nameStr != "" {
		tf.Name = nameStr
	}

	tagsStr := query["tags[]"]
	if len(tagsStr) > 0 {

		tags := make([]string, 0, len(tagsStr))
		for _, tag := range tagsStr {
			_, err := uuid.Parse(tag)
			if err != nil {
				return fmt.Errorf("invalid tag id: %s", tag)
			}
			tags = append(tags, tag)

		}
		tf.Tags = tags
	}

	isAsc, err := strconv.ParseBool(query.Get("isAsc"))
	if err != nil {
		isAsc = true
	}
	tf.IsAsc = isAsc

	sortBy := query.Get("sortBy")
	if sortBy != "" {
		filter := SortByFilter(sortBy)
		err := filter.Parse(sortBy)
		if err != nil {
			return fmt.Errorf("invalid sortBy value: %s", sortBy)
		}
		tf.SortBy = filter
	}

	priceStr := query["price[]"]
	if len(priceStr) > 0 {
		if len(priceStr) != 2 {
			if err != nil {
				return fmt.Errorf("invalid price format. Expected 2 values")
			}
		}
		prices := make([]float64, 0, len(priceStr))
		for _, p := range priceStr {
			parsedPrice, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return fmt.Errorf("invalid price: %s", p)
			}
			prices = append(prices, parsedPrice)
			sort.Float64s(prices)
		}
		tf.MinPrice = prices[0]
		tf.MaxPrice = prices[1]
	}
	return nil
}
