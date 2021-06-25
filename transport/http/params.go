package http

import (
	"net/http"
	"strconv"

	"github.com/Financial-Times/gourmet/apperror"
	"github.com/gorilla/mux"
)

func ParseIntPathParam(req *http.Request, paramName string, paramDesc string) (int, error) {
	vars := mux.Vars(req)
	id, ok := vars[paramName]
	if !ok {
		return 0, apperror.ValidationError.Newf("missing or invalid %s %s", paramDesc, id)
	}
	p, _ := strconv.ParseInt(id, 10, 64)

	return int(p), nil
}

func ParseUintQueryParam(req *http.Request, paramName string) uint {
	q := req.URL.Query()
	p, _ := strconv.ParseUint(q.Get(paramName), 10, 64)
	return uint(p)
}
