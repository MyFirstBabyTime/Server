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