package grpcgw

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang-tire/pkg/log"
)

type swaggerFile struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`

	Host        string                 `json:"host"`
	Schemes     []string               `json:"schemes"`
	Consumes    []string               `json:"consumes"`
	Produces    []string               `json:"produces"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
}

type swaggerServer struct {
	swaggerBaseURL string
}

var (
	swaggerLock sync.RWMutex
	data        = swaggerFile{
		Swagger: "2.0",
		Info: struct {
			Title   string `json:"title"`
			Version string `json:"version"`
		}{Title: "Swagger", Version: "1.0"},
		Schemes:     []string{"https", "http"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Host:        "",
		Paths:       make(map[string]interface{}),
		Definitions: make(map[string]interface{}),
	}
)

func RegisterSwagger(paths, definitions map[string]interface{}) {
	swaggerLock.Lock()
	defer swaggerLock.Unlock()

	for k, v := range paths {
		data.Paths[k] = v
	}

	for k, v := range definitions {
		data.Definitions[k] = v
	}
}

func (sw *swaggerServer)swaggerHandler(w http.ResponseWriter, r *http.Request) {
	swaggerLock.Lock()
	defer swaggerLock.Unlock()

	reqFile := strings.TrimPrefix(r.RequestURI, sw.swaggerBaseURL)
	if reqFile == "index.json" {
		data.Host = r.Host
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error("serve swagger json failed", log.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if reqFile == "" {
		reqFile = "index.html"
	}

	fileData, err := Asset(reqFile)
	if err != nil {
		log.Error("asset file not found", log.Err(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch strings.ToLower(filepath.Ext(reqFile)) {
	case ".json":
		w.Header().Add("Content-Type", "application/json")
	case ".css":
		w.Header().Add("Content-Type", "text/css")
	case ".js":
		w.Header().Add("Content-Type", "text/javascript")
	case ".html":
		w.Header().Add("Content-Type", "text/html")
	}
	_, _ = w.Write(fileData)
}
