package web

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	PageDefault     = 0
	PageSizeDefault = 50
	PageSizeMaximum = 200
)

type Params struct {
	UserId            string
	TokenId           string
	PartId            string
	RoleId            string
	RoleName          string
	Scope             string
	Page              int
	PageSize          int
	IncludeTotalCount bool
	Q                 string
}

func ParamsFromRequest(r *http.Request) *Params {
	params := &Params{}

	props := mux.Vars(r)
	query := r.URL.Query()

	params.UserId = props["user_id"]
	params.TokenId = props["token_id"]
	params.PartId = props["part_id"]
	params.RoleId = props["role_id"]
	params.RoleName = props["role_name"]
	params.Scope = query.Get("scope")

	if val, err := strconv.Atoi(query.Get("page")); err != nil || val < 0 {
		params.Page = PageDefault
	} else {
		params.Page = val
	}

	params.PageSize = getPageSizeFromQuery(query)

	params.IncludeTotalCount, _ = strconv.ParseBool(query.Get("include_total_count"))

	params.Q = query.Get("q")

	return params
}

func getPageSizeFromQuery(query url.Values) int {
	val, err := strconv.Atoi(query.Get("per_page"))
	if err != nil {
		val, err = strconv.Atoi(query.Get("pageSize"))
	}

	if err != nil || val < 0 {
		return PageSizeDefault
	} else if val > PageSizeMaximum {
		return PageSizeMaximum
	}

	return val
}
