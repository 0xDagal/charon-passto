// Package passto (an passto plugin).
package passto

import (
    "context"
    "net/http"
    "encoding/json"
    "strings"
    "github.com/elastic/go-elasticsearch/v8"
)

// Config the plugin configuration.
type Config struct { }

// passto a plugin.
type Passto struct {
    next     http.Handler
    name     string
    es       elasticsearch.Client
}

type Log struct {
    Method      string
    RemoteAddr  string
    RequestURI  string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config { return &Config{ } }

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
    cfg := elasticsearch.Config{
        Addresses: []string{
            "http://elasticsearch-master:9200",
        },
        Password: "changeme",
    }
    client, _ := elasticsearch.NewClient(cfg)
    return &Passto{
        next: next,
		name: name,
        es  : *client,
    }, nil
}

func (p *Passto) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    rw.Write([]byte(req.RemoteAddr))
    log := Log{
        Method: req.Method,
        RemoteAddr: req.RemoteAddr,
        RequestURI: req.RequestURI,
    }
    json, _ := json.Marshal(log)
    res, _ := p.es.Index(
        "log",
        strings.NewReader(string(json)),
        p.es.Index.WithDocumentID("1"),
        p.es.Index.WithRefresh("true"),
    )
    defer res.Body.Close()
}
