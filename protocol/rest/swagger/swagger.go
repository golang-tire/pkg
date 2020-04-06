// Package swagger this package come from fzerorubigd engine package github.com/fzerorubigd/engine
// with some changes
package swagger

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang-tire/pkg/log"
	"github.com/golang-tire/pkg/types"
	"go.uber.org/zap"
)

type swaggerFile struct {
	Swagger string      `json:"swagger"`
	Info    types.MapSI `json:"info"`

	Host        string      `json:"host"`
	Schemes     []string    `json:"schemes"`
	Consumes    []string    `json:"consumes"`
	Produces    []string    `json:"produces"`
	Paths       types.MapSI `json:"paths"`
	Definitions types.MapSI `json:"definitions"`
}

var (
	swaggerLock sync.RWMutex
	data        = swaggerFile{
		Swagger:     "2.0",
		Info:        make(types.MapSI),
		Schemes:     []string{"https", "http"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Host:        "",
		Paths:       make(types.MapSI),
		Definitions: make(types.MapSI),
	}
)

// Server swagger server interface
type Server interface {
	Handler(w http.ResponseWriter, r *http.Request)
	Register(paths types.MapSI, definitions types.MapSI, info types.MapSI, userInfo bool)
}

type server struct {
	logger log.Logger
}

// New create and return a server object
func New(logger log.Logger) Server {
	return &server{logger}
}

// Handler register swagger handler
func (s server) Handler(w http.ResponseWriter, r *http.Request) {

	swaggerLock.RLock()
	defer swaggerLock.RUnlock()

	fl := strings.TrimPrefix(r.RequestURI, "/v1/swagger/")
	s.logger.Info("Swagger file requested", zap.String("file", fl))
	if fl == "index.json" {
		w.Header().Add("Content-Type", "application/json")
		var d = data
		d.Host = r.Host
		if err := json.NewEncoder(w).Encode(d); err != nil {
			s.logger.Error("Failed to serve swagger", zap.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if fl == "" {
		fl = "index.html"
	}

	d, err := Asset(fl)
	if err != nil {
		s.logger.Error("Not found", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusNotFound)
	}
	switch strings.ToLower(filepath.Ext(fl)) {
	case ".json":
		w.Header().Add("Content-Type", "application/json")
	case ".css":
		w.Header().Add("Content-Type", "text/css")
	case ".js":
		w.Header().Add("Content-Type", "text/javascript")
	case ".html":
		w.Header().Add("Content-Type", "text/html")
	}
	_, _ = w.Write(d)
}

// Register, register a swagger end point
func (s server) Register(paths types.MapSI, definitions types.MapSI, info types.MapSI, userInfo bool) {
	swaggerLock.Lock()
	defer swaggerLock.Unlock()

	for i := range paths {
		_, ok := data.Paths[i]
		if !ok {
			s.logger.Error("Path is already registered")
		}
		data.Paths[i] = appendSecurity(paths[i], strings.Contains(i, "{"))
	}

	for i := range definitions {
		_, ok := data.Definitions[i]
		if !ok {
			s.logger.Error("Definition is already registered")
		}
		data.Definitions[i] = definitions[i]
	}

	if userInfo {
		data.Info = info
	}

	if data.Definitions["ErrorResponse"] == nil {
		data.Definitions["ErrorResponse"] = errDef
		data.Definitions["ErrorFldResponse"] = errFldDef
	}
}

func appendSecurity(d interface{}, has404 bool) types.MapSI {
	v := d.(types.MapSI)
	for i := range v {
		meth, ok := v[i].(types.MapSI)
		if !ok {
			continue
		}
		if meth["security"] != nil {
			if p, ok := meth["parameters"].([]interface{}); !ok {
				meth["parameters"] = createParameter(nil)
			} else {
				meth["parameters"] = createParameter(p)
			}
		}

		meth["responses"] = create40XResponses(meth["responses"].(types.MapSI), meth["security"] != nil, has404)
	}
	return v
}

func createParameter(old []interface{}) []interface{} {
	return append(old, types.MapSI{
		"description": "the security token, get it from login route",
		"in":          "header",
		"name":        "authorization",
		"required":    true,
		"type":        "string",
	})
}

func create40XResponses(in types.MapSI, forbidden, notFound bool) types.MapSI {
	if forbidden {
		in["401"] = err401
		in["403"] = err403
	}
	if notFound {
		in["404"] = err404
	}
	in["400"] = err400
	return in
}
