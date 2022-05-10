package passto

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Config the plugin configuration.
type Config struct {
	ES elasticsearch.Config
}

// passto a plugin.
type Passto struct {
	next http.Handler
	name string
	es   *elasticsearch.Client
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	return &Config{
		ES: cfg,
	}
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	client, err := elasticsearch.NewClient(config.ES)
	if err != nil {
		log.Panicf("Could not create client: %s\n", err)
	}
	return &Passto{
		next: next,
		name: name,
		es:   client,
	}, nil
}

func (p *Passto) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	type Log struct {
		Method        string
		Proto         string
		ContentLength int64
		Host          string
		RemoteAddr    string
		RequestURI    string
	}
	reqlog := Log{
		Method:        req.Method,
		Proto:         req.Proto,
		ContentLength: req.ContentLength,
		Host:          req.Host,
		RemoteAddr:    req.RemoteAddr,
		RequestURI:    req.RequestURI,
	}
	body, err := json.Marshal(reqlog)
	if err != nil {
		log.Panicf("Could not marshal body: %s\n", err)
	}
	esreq := esapi.IndexRequest{
		Index:      "test",
		DocumentID: "1",
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := esreq.Do(context.Background(), p.es)
	if err != nil {
		log.Panicf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Panicf("Error parsing the response body: %s", err)
		} else {
			log.Panicf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	p.next.ServeHTTP(rw, req)
}
