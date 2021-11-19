package hh

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/mcsolutions/lib/back/common/header"
	"gitlab.com/mcsolutions/lib/back/common/json"
	"gitlab.com/mcsolutions/lib/back/common/logger"
	"gitlab.com/mcsolutions/lib/back/common/status"
	"io"
	"net/http"
	"strings"
)

type (
	Route struct {
		Pattern     string
		Methods     string
		HandlerFunc http.HandlerFunc
	}
)

const (
	OPTIONS = "," + http.MethodOptions
)

func NewRouter(routes []Route) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc
		methods := strings.Split(route.Methods+OPTIONS, ",")
		router.
			Methods(methods...).
			Path(route.Pattern).
			Handler(handler)
	}
	return router
}

func GetHandlerFunc(f func(*http.Request) error, logger *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header.SetCorsType(&w)
		if r.Method == http.MethodOptions {
			header.SetOtherOptionsHeader(&w)
			w.WriteHeader(http.StatusNoContent)
		} else if err := f(r); err != nil {
			header.SetContentTypeJson(&w)
			w.WriteHeader(http.StatusBadRequest)
			status.WriteLoggerStatus(&w, err, logger)
		} else {
			header.SetContentTypeJson(&w)
			status.WriteStatus(&w, nil)
		}
	}
}

func GetHandlerWFunc(f func(http.ResponseWriter, *http.Request) error, logger *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header.SetCorsType(&w)
		if r.Method == http.MethodOptions {
			header.SetOtherOptionsHeader(&w)
			w.WriteHeader(http.StatusNoContent)
		} else if err := f(w, r); err != nil {
			header.SetContentTypeJson(&w)
			status.WriteLoggerStatus(&w, err, logger)
		} else {
			header.SetContentTypeJson(&w)
			status.WriteStatus(&w, nil)
		}
	}
}

func GetCodeHandlerFunc(f func(*http.Request) (int, error), logger *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header.SetCorsType(&w)
		if r.Method == http.MethodOptions {
			header.SetOtherOptionsHeader(&w)
			w.WriteHeader(http.StatusNoContent)
		} else if code, err := f(r); err != nil {
			header.SetContentTypeJson(&w)
			w.WriteHeader(code)
			status.WriteLoggerStatus(&w, err, logger)
		} else {
			header.SetContentTypeJson(&w)
			w.WriteHeader(code)
			status.WriteStatus(&w, nil)
		}
	}
}

func GetObjHandlerFunc(f func(*http.Request) (interface{}, error), logger *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header.SetCorsType(&w)
		if r.Method == http.MethodOptions {
			header.SetOtherOptionsHeader(&w)
			w.WriteHeader(http.StatusNoContent)
		} else if obj, err := f(r); err != nil {
			header.SetContentTypeJson(&w)
			w.WriteHeader(http.StatusBadRequest)
			status.WriteLoggerStatus(&w, err, logger)
		} else if obj == nil {
			w.WriteHeader(http.StatusNoContent)
		} else if bytes, ok := obj.([]byte); ok {
			w.Write(bytes)
		} else {
			header.SetContentTypeJson(&w)
			if err := json.Encode(obj, io.Writer(w)); err != nil {
				status.StatusJsonEncodeError(&w, err, logger)
			}
		}
	}
}

func GetCodeObjHandlerFunc(f func(*http.Request) (int, interface{}, error), logger *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header.SetCorsType(&w)
		if r.Method == http.MethodOptions {
			header.SetOtherOptionsHeader(&w)
			w.WriteHeader(http.StatusNoContent)
		} else if code, obj, err := f(r); err != nil {
			header.SetContentTypeJson(&w)
			w.WriteHeader(code)
			status.WriteLoggerStatus(&w, err, logger)
		} else if obj == nil {
			if code == http.StatusOK {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(code)
			}
		} else if bytes, ok := obj.([]byte); ok {
			w.WriteHeader(code)
			w.Write(bytes)
		} else {
			header.SetContentTypeJson(&w)
			w.WriteHeader(code)
			if err := json.Encode(obj, io.Writer(w)); err != nil {
				status.StatusJsonEncodeError(&w, err, logger)
			}
		}
	}
}

func GetCodeObjFileHandlerFunc(f func(*http.Request) (int, interface{}, string, error), logger *logger.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		header.SetCorsType(&w)
		if r.Method == http.MethodOptions {
			header.SetOtherOptionsHeader(&w)
			w.WriteHeader(http.StatusNoContent)
		} else if code, obj, filename, err := f(r); err != nil {
			header.SetContentTypeJson(&w)
			w.WriteHeader(code)
			status.WriteLoggerStatus(&w, err, logger)
		} else if obj == nil {
			if code == http.StatusOK {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(code)
			}
		} else if bytes, ok := obj.([]byte); ok {
			w.WriteHeader(code)
			w.Write(bytes)
		} else {
			w.Header().Set("Content-Disposition",
				fmt.Sprintf("attachment; filename=%s", filename))
			http.ServeFile(w, r, filename)
		}
	}
}
