package request

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET_PARAM  string = "GET"
	POST_PARAM string = "POST"
	PATH_PARAM string = "PATH"
)

func GetParameterValue(r *http.Request, param_type, param_name string) string {

	switch param_type {
	case GET_PARAM:
		return r.URL.Query().Get(param_name)
	case POST_PARAM:
		return r.PostFormValue(param_name)
	case PATH_PARAM:
		path_vars := mux.Vars(r)
		return path_vars[param_name]
	default:
		return ""
	}
}

func GetParameterValues(r *http.Request, param_type string) []string {

	switch param_type {
	case GET_PARAM:
		return r.URL.Query()[param_type]
	case POST_PARAM:
		return r.PostForm[param_type]
	default:
		return nil
	}
}
