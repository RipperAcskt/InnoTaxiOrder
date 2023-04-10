package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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

func (es *Elastic) CreateOrder(ctx context.Context, order model.Order) (string, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := json.Marshal(order)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %s", err)
	}
	req := esapi.IndexRequest{
		Index:   es.cfg.ELASTIC_DB_NAME,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(queryCtx, es.Client)
	if err != nil {
		return "", fmt.Errorf("req do failed: %w", err)
	}
	defer res.Body.Close()

	var idStruct struct {
		Id string `json:"_id"`
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read all failed: %w", err)
	}
	err = json.Unmarshal(body, &idStruct)
	if err != nil {
		return "", fmt.Errorf("unmarshal failed: %w", err)
	}

	if res.IsError() {
		return "", fmt.Errorf("res error: %w", err)
	}

	return idStruct.Id, nil
}

func (es *Elastic) GetOrders(ctx context.Context, indexes []string) ([]*model.Order, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if len(indexes) == 0 {
		res, err := es.Client.Search(
			es.Client.Search.WithContext(queryCtx),
			es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
		)
		if err != nil {
			return nil, fmt.Errorf("search failed: %w", err)
		}
		defer res.Body.Close()
		return es.parseInfo(res)
	}

	var body bytes.Buffer
	query := map[string]interface{}{
		"ids": map[string]interface{}{
			"values": indexes,
		},
	}
	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return nil, fmt.Errorf("encode failed: %w", err)
	}
	b, _ := io.ReadAll(&body)
	fmt.Println(string(b))
	res, err := es.Client.Search(
		es.Client.Search.WithContext(queryCtx),
		es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
		es.Client.Search.WithBody(&body),
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()
	return es.parseInfo(res)
}

func (es *Elastic) parseInfo(res *esapi.Response) ([]*model.Order, error) {

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("decode failed: %w", err)
		} else if err != nil {
			return nil, fmt.Errorf("error: %w", err)
		}
	}

	var info model.ElasticModel
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	var orders []*model.Order
	for _, el := range info.Hits.Hits {
		element := el
		element.Source.ID = element.ID
		orders = append(orders, &element.Source)
	}
	return orders, nil
}
