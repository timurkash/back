package httpwalker

import (
	"github.com/gorilla/mux"
	"gitlab.com/mcsolutions/lib/back/common/others"
	"net/http"
	"strings"
)

type (
	Route struct {
		Name        string
		Methods     string
		PathPrefix  string
		MSUrl       string
		HandlerFunc http.HandlerFunc
	}
	Additional struct {
		Name       string
		PathPrefix string
		MSUrl      string
	}
)

var additionals []Additional

func NewRouterPrefix(routes []Route) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		//handler = logger.RestLogger(handler, route.Name)
		router.
			Methods(strings.Split(route.Methods, ",")...).
			PathPrefix(route.PathPrefix).
			Name(route.Name).
			Handler(handler)
		additionals = append(additionals, Additional{route.Name, route.PathPrefix, route.MSUrl})
	}
	return router
}

func GetRoute(pattern, method string, routes []Route) *Route {
	found := false
	result := Route{}
	for _, route := range routes {
		if route.PathPrefix == pattern && others.Contains(strings.Split(route.Methods, ","), method) {
			found = true
			result = route
		}
	}
	if !found {
		return nil
	}
	return &result
}

func GetAdditional(name string) *Additional {
	for _, additional := range additionals {
		if additional.Name == name {
			return &additional
		}
	}
	return nil
}

func GetVars(r *http.Request, muxVar string) (value string) {
	vars := mux.Vars(r)
	value = vars[muxVar]
	return
}
