package hh

import (
	"gitlab.com/mcsolutions/lib/back/common/env"
	"gitlab.com/mcsolutions/lib/back/common/header"
	"gitlab.com/mcsolutions/lib/back/common/json"
	"gitlab.com/mcsolutions/lib/back/common/logger"
	"gitlab.com/mcsolutions/lib/back/common/status"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type DevOpsRouters struct {
	Started     *time.Time
	SubgroupUrl string
	Name        string
	Version     string
	Revision    string
	Envs        []env.Env
	Logger      *logger.Logger
	IsReady     *atomic.Value
}

func (d *DevOpsRouters) healthCheckRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/healthz",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		},
	}
}

func (d *DevOpsRouters) readyCheckRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/readyz",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			if d.IsReady == nil || !d.IsReady.Load().(bool) {
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		},
	}
}

func (d *DevOpsRouters) gitlabRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/gitlab",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			url := d.SubgroupUrl + d.Name
			bytes := []byte("<a target='_blank' href='" + url + "'>" + url + "</a>")
			w.Write(bytes)
		},
	}
}

type RuntimeResponse struct {
	Started      *time.Time `json:"started"`
	Hostname     string     `json:"hostname,omitempty"`
	Version      string     `json:"version,omitempty"`
	Revision     string     `json:"revision,omitempty"`
	Message      string     `json:"message,omitempty"`
	GOOS         string     `json:"goos,omitempty"`
	GOARCH       string     `json:"goarch,omitempty"`
	Runtime      string     `json:"runtime,omitempty"`
	NumGoroutine string     `json:"num_goroutine,omitempty"`
	NumCPU       string     `json:"num_cpu,omitempty"`
}

func (d *DevOpsRouters) infoRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			header.SetContentTypeJson(&w)
			data := &RuntimeResponse{
				Started:      d.Started,
				Version:      d.Version,
				Revision:     d.Revision,
				GOOS:         runtime.GOOS,
				GOARCH:       runtime.GOARCH,
				Runtime:      runtime.Version(),
				NumGoroutine: strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
				NumCPU:       strconv.FormatInt(int64(runtime.NumCPU()), 10),
			}
			if err := json.Encode(data, w); err != nil {
				status.StatusJsonEncodeError(&w, err, d.Logger)
			}
		},
	}
}

func (d *DevOpsRouters) envsRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/envs",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			header.SetContentTypeJson(&w)
			env := os.Environ()
			if err := json.Encode(env, w); err != nil {
				status.StatusJsonEncodeError(&w, err, d.Logger)
			}
		},
	}
}

func (d *DevOpsRouters) myEnvsRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/myenvs",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			if err := yaml.NewEncoder(w).Encode(d.Envs); err != nil {
				status.StatusJsonEncodeError(&w, err, d.Logger)
			}
		},
	}
}

func (d *DevOpsRouters) myEnvsJsonRoute() Route {
	return Route{
		Methods: "GET",
		Pattern: "/myenvs/json",
		HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
			if len(d.Envs) > 0 {
				header.SetContentTypeJson(&w)
				if err := yaml.NewEncoder(w).Encode(d.Envs); err != nil {
					status.StatusJsonEncodeError(&w, err, d.Logger)
				}
			}
		},
	}
}

func (d *DevOpsRouters) AddTo(routes []Route) []Route {
	result := append(routes,
		d.gitlabRoute(),
		d.healthCheckRoute(),
		d.readyCheckRoute(),
		d.infoRoute(),
		d.envsRoute(),
		d.myEnvsRoute(),
		d.myEnvsJsonRoute())
	return result
}
