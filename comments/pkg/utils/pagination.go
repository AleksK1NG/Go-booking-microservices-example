package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultSize = 10
)

// Pagination query params
type Pagination struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

// NewPaginationQuery
func NewPaginationQuery(size int, page int) *Pagination {
	return &Pagination{Size: size, Page: page}
}

// Set page size
func (q *Pagination) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

// Set page number
func (q *Pagination) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// Set order by
func (q *Pagination) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// Get offset
func (q *Pagination) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// Get limit
func (q *Pagination) GetLimit() int {
	return q.Size
}

// Get OrderBy
func (q *Pagination) GetOrderBy() string {
	return q.OrderBy
}

// Get OrderBy
func (q *Pagination) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *Pagination) GetSize() int {
	return q.Size
}

func (q *Pagination) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

// Get pagination query struct from
func GetPaginationFromCtx(c echo.Context) (*Pagination, error) {
	q := &Pagination{}
	if err := q.SetPage(c.QueryParam("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.QueryParam("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.QueryParam("orderBy"))

	return q, nil
}

// Get total pages int
func (q *Pagination) GetTotalPages(totalCount int) int {
	// d := float64(totalCount) / float64(pageSize)
	d := float64(totalCount) / float64(q.GetSize())
	return int(math.Ceil(d))
}

// Get has more
func (q *Pagination) GetHasMore(totalCount int) bool {
	return q.GetPage() < totalCount/q.GetSize()
}
