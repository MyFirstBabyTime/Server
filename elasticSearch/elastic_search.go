package elasticSearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type elasticSearch struct {
	es *elasticsearch.Client
}

func New(esEndPoint string) *elasticSearch {
	cfg := elasticsearch.Config{
		Addresses: []string{esEndPoint},
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect elasticSearch").Error())
	}

	return &elasticSearch{
		es,
	}
}
