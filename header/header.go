package header

import (
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/mcsolutions/lib/back/common/consts"
	"net/http"
)

func SetContentTypeJson(w *http.ResponseWriter) {
	(*w).Header().Set(consts.CONTENT_TYPE, consts.CONTENT_TYPE_JSON)
}

func SetCorsType(w *http.ResponseWriter) {
	(*w).Header().Set(consts.CORS_TYPE, "*")
}

func SetOtherOptionsHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "*") //Origin, Content-Type, X-Auth-Token
}

func getErrorNotSpecified(name string) error {
	return errors.New(name + " not specified")
}

func GetHeaderRequired(r *http.Request, name string) (string, error) {
	value := r.Header.Get(name)
	if value == "" {
		return "", getErrorNotSpecified(name)
	}
	return value, nil
}

func GetVar(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}

func GetVarRequired(r *http.Request, name string) (string, error) {
	value := GetVar(r, name)
	if value == "" {
		return "", getErrorNotSpecified(name)
	}
	return value, nil
}

func Get2Vars(r *http.Request, name1, name2 string) (string, string) {
	vars := mux.Vars(r)
	return vars[name1], vars[name2]
}

func GetVars(r *http.Request, names ...string) []string {
	vars := mux.Vars(r)
	result := make([]string, 0)
	for _, name := range names {
		result = append(result, vars[name])
	}
	return result
}
