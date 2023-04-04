package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type Elastic struct {
	Client *elasticsearch.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Elastic, error) {
	cfgEs := elasticsearch.Config{
		Addresses: []string{cfg.ELASTIC_DB_HOST},
		Username:  cfg.ELASTIC_DB_USERNAME,
		Password:  cfg.ELASTIC_DB_PASSWORD,
	}
	es, err := elasticsearch.NewClient(cfgEs)
	if err != nil {
		return nil, fmt.Errorf("new client failed: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("info failed: %s", err)
	}

	defer res.Body.Close()
	return &Elastic{es, cfg}, nil
}

func (es *Elastic) CreateOrder(ctx context.Context, order model.Order) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal failed: %s", err)
	}
	req := esapi.IndexRequest{
		Index:   es.cfg.ELASTIC_DB_NAME,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(queryCtx, es.Client)
	if err != nil {
		return fmt.Errorf("req do failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("res error: %w", err)
	}

	return nil
}

// func (es *Elastic) GetOrders(ctx context.Context) (*model.ElasticModel, error) {
// 	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	res, err := es.Client.Search(
// 		es.Client.Search.WithContext(queryCtx),
// 		es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("search failed: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		var e map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
// 			return nil, fmt.Errorf("decode failed: %w", err)
// 		} else {
// 			return nil, fmt.Errorf("error: %w", err)
// 		}
// 	}

// 	s, _ := io.ReadAll(res.Body)
// 	fmt.Println(string(s))
// 	var info model.ElasticModel
// 	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
// 		return nil, fmt.Errorf("decode failed: %w", err)
// 	}

// 	return &info, nil
// }
