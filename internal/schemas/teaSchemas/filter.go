package teaSchemas

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sort"
	"strconv"
)

type Filters struct {
	Limit         uint64       `json:"limit" db:"limit"`
	Page          uint64       `json:"page"`
	Offset        uint64       `db:"offset"`
	CategoryId    uuid.UUID    `json:"categoryId,omitempty" db:"category_id"`
	Name          string       `json:"name,omitempty" db:"name"`
	Tags          []string     `json:"tags,omitempty" db:"tags"`
	MinServePrice float64      `json:"minServePrice,omitempty" db:"min_serve_price"`
	MaxServePrice float64      `json:"maxServePrice,omitempty" db:"max_serve_price"`
	SortBy        SortByFilter `json:"sortBy,omitempty"`
	IsAsc         bool         `json:"isAsc"`
	IsDeleted     bool         `json:"isDeleted,omitempty" db:"is_deleted"`
	UserId        uuid.UUID    `db:"user_id"`
}

type SortByFilter string

const (
	Name       SortByFilter = "name"
	ServePrice SortByFilter = "servePrice"
	Rating     SortByFilter = "rating"
)

func (f *SortByFilter) String() string {
	return string(*f)
}
func (f *SortByFilter) Parse(s string) error {
	SortByMapping := map[string]SortByFilter{
		"name":       Name,
		"servePrice": ServePrice,
		"rating":     Rating,
	}
	if val, ok := SortByMapping[s]; ok {
		*f = val
		return nil
	}
	return fmt.Errorf("invalid TeaSortByFilter value: %s", s)
}
func (f *SortByFilter) ToDbFilter() string {
	dbFilters := map[SortByFilter]string{
		Name:       "name",
		ServePrice: "serve_price",
		Rating:     "rating",
	}
	return dbFilters[*f]
}

func (tf *Filters) Validate(r *http.Request) error {
	query := r.URL.Query()
	limit, err := strconv.ParseUint(query.Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	page, err := strconv.ParseUint(query.Get("page"), 10, 64)
	if err != nil {
		page = 1
	}

	if page == 0 {
		return fmt.Errorf("the page can not be equal to 0")
	}

	tf.Limit = limit
	tf.Page = page

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
		mappedSortBy := SortByFilter(sortBy)
		err := mappedSortBy.Parse(sortBy)
		if err != nil {
			return fmt.Errorf("invalid sortBy value: %s", sortBy)
		}
		tf.SortBy = mappedSortBy
	}

	servePriceStr := query["servePrice[]"]
	if len(servePriceStr) > 0 {
		if len(servePriceStr) != 2 {
			if err != nil {
				return fmt.Errorf("invalid servePrice format. Expected 2 values")
			}
		}
		prices := make([]float64, 0, len(servePriceStr))
		for _, p := range servePriceStr {
			parsedPrice, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return fmt.Errorf("invalid servePrice: %s", p)
			}
			prices = append(prices, parsedPrice)
			sort.Float64s(prices)
		}
		tf.MinServePrice = prices[0]
		tf.MaxServePrice = prices[1]
	}

	isDeleted, err := strconv.ParseBool(query.Get("isDeleted"))
	if err != nil {
		isDeleted = false
	}
	tf.IsDeleted = isDeleted

	return nil
}
