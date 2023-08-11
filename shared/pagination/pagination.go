package pagination

import (
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/evermos/boilerplate-go/shared/failure"
)

func ConvertToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

type Pagination struct {
	Page   int    `validate:"required"`
	Limit  int    `validate:"required"`
	Offset int    `db:"offset"`
	Field  string `db:"field"`
	Sort   string `db:"sort"`
}

func NewPaginationQuery(page, limit int, field, sort string) *Pagination {
	pg := Pagination{Page: page, Limit: limit, Offset: (page - 1) * limit, Field: field, Sort: sort}
	return &pg
}

func GetSortDirection(s string) string {
	switch strings.ToLower(s) {
	case "asc":
		return "ASC"
	case "desc":
		return "DESC"
	default:
		return "ASC"
	}
}

func CheckFieldQuery(s string, def string) string {
	switch strings.ToLower(s) {
	case "":
		return strings.ToLower(def)
	default:
		return strings.ToLower(s)
	}
}

func ParseQueryParams(r *http.Request, key string) string {
	return strings.ToLower(r.URL.Query().Get(key))
}

func GetCancelled(str string) (booleanValue bool) {
	switch strings.ToLower(str) {
	case "true":
		booleanValue = true
		return
	case "false":
		booleanValue = false
		return
	default:
		booleanValue = false
		return
	}
}

func GetPagination(r *http.Request) (pg *Pagination, err error) {
	page, err := ConvertToInt(ParseQueryParams(r, "page"))
	if err != nil {
		err = failure.BadRequest(err)
		return
	}
	limit, err := ConvertToInt(ParseQueryParams(r, "limit"))
	if err != nil {
		err = failure.BadRequest(err)
		return
	}
	sort := GetSortDirection(ParseQueryParams(r, "sort"))
	field := CheckFieldQuery(ParseQueryParams(r, "field"), "id")
	pg = NewPaginationQuery(page, limit, field, sort)
	return
}

func (p *Pagination) GetTotalPages(res interface{}) int {
	val := reflect.ValueOf(res)
	if val.Kind() != reflect.Slice {
		return 0
	}
	length := val.Len()
	return int(math.Ceil(float64(length) / float64(p.Limit)))
}
